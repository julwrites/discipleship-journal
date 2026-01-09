package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"discipleship_journal_api/database"
	"discipleship_journal_api/handlers"
	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"

	_ "discipleship_journal_api/docs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Discipleship Journal API
// @version 1.0
// @description API for the Discipleship Journal application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// mockDatabase implements services.DBInterface but returns errors for all operations
// Used when database connection fails
type mockDatabase struct {
	connected bool
}

func (m *mockDatabase) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, fmt.Errorf("database not connected")
}

func (m *mockDatabase) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	// Return a mock row that will error on Scan
	return &mockRow{}
}

func (m *mockDatabase) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, fmt.Errorf("database not connected")
}

func (m *mockDatabase) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, fmt.Errorf("database not connected")
}

// mockRow implements pgx.Row but returns error on Scan
type mockRow struct{}

func (m *mockRow) Scan(dest ...any) error {
	return fmt.Errorf("database not connected")
}

// @host localhost:8080
// @BasePath /
func main() {
	// Load .env file if it exists (local dev)
	_ = godotenv.Load()

	// Setup Logger
	var logger *slog.Logger
	if os.Getenv("APP_ENV") == "production" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	// Initialize Secret Loader
	secretLoader, err := services.NewSecretLoader(context.Background(), "")
	if err != nil {
		logger.Error("Failed to initialize secret loader", "error", err)
		// Continue anyway - secret loader will fall back to env vars
	}
	if secretLoader != nil {
		defer secretLoader.Close()
	}

	// Load configuration from Secret Manager or environment variables
	loadSecret := func(secretName string) string {
		var value string
		if secretLoader != nil {
			value, _ = secretLoader.LoadSecret(context.Background(), secretName)
			if value == "" {
				value = os.Getenv(secretName)
			}
		} else {
			value = os.Getenv(secretName)
		}
		return value
	}

	// Load individual database components
	dbSecrets := []string{"DB_USERNAME", "DB_PASSWORD", "DB_NAME", "DB_HOST", "DB_PORT", "CLOUD_SQL_INSTANCE"}
	for _, secret := range dbSecrets {
		value := loadSecret(secret)
		if value != "" {
			os.Setenv(secret, value)
			// Mask password in logs
			if secret == "DB_PASSWORD" {
				logger.Info("Loaded secret", "secret", secret, "value", "***")
			} else {
				logger.Info("Loaded secret", "secret", secret, "value", value)
			}
		}
	}

	// Load CORS allowed origins from Secret Manager
	corsOrigins := loadSecret("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		os.Setenv("CORS_ALLOWED_ORIGINS", corsOrigins)
		logger.Info("Loaded CORS allowed origins", "value", corsOrigins)
	}

	// Try to connect to database, but don't exit immediately in production
	// Cloud Run needs the container to start listening on the port first
	dbConnected := false
	if err := database.Connect(); err != nil {
		logger.Error("Database connection failed", "error", err)
		// Don't exit immediately - let the server start and health check will fail
		// This allows Cloud Run to properly start the container
		if os.Getenv("APP_ENV") == "production" {
			logger.Warn("Database connection failed in production, but continuing to start server")
		}
	} else {
		dbConnected = true
		defer database.Close()

		// Run database migrations if connected
		if err := database.RunMigrations(); err != nil {
			logger.Error("Database migrations failed", "error", err)
			// In production, we might want to exit if migrations fail as schema might be incompatible.
			// However, for robustness, we'll log and continue, unless it's a critical error.
			// The user can check logs.
		}
	}

	// Init Firebase Service
	// Pass empty string for saKey to use default credentials (production)
	// or rely on GOOGLE_APPLICATION_CREDENTIALS
	firebaseProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if firebaseProjectID == "" {
		firebaseProjectID = "mock-project-id" // Use environment variable GOOGLE_CLOUD_PROJECT for real project ID
	}
	firebaseService, err := services.NewFirebaseService(context.Background(), "", firebaseProjectID)
	if err != nil {
		logger.Error("Firebase init failed", "error", err)
	}

	var authMiddleware func(http.Handler) http.Handler
	var notificationService services.NotificationService

	if firebaseService != nil {
		authMiddleware = middleware.NewAuthMiddlewareFromClient(firebaseService.AuthClient).VerifyToken
		if firebaseService.MessagingClient != nil {
			notificationService = services.NewNotificationService(database.DB, firebaseService.MessagingClient)
		} else {
			logger.Warn("Firebase Messaging not initialized")
			notificationService = services.NewMockNotificationService() // Fallback to avoid nil pointer
		}
	} else {
		// Fallback for local dev without firebase creds
		// Create a mock auth middleware that injects a test user
		logger.Warn("Firebase Service is nil, using mock authentication for development")
		authMiddleware = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For development, accept any request (or require simple test token)
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" || authHeader == "Bearer test" {
					// Inject a test user ID into context using TestUserKey
					// This matches what GetUserUUID expects for test mode
					ctx := context.WithValue(r.Context(), handlers.TestUserKey, "00000000-0000-0000-0000-000000000001")
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					http.Error(w, "Unauthorized - use 'Bearer test' for development mode", http.StatusUnauthorized)
				}
			})
		}
		// Initialize mock notification service to avoid nil pointer panics
		notificationService = services.NewMockNotificationService()
	}

	// Init Bible AI Client
	var bibleAIClient services.BibleAIClient

	// Load Bible API secrets from Secret Manager or environment variables
	var bibleAPIURL, bibleAPIKey, llmSystemPrompts string
	if secretLoader != nil {
		bibleAPIURL, _ = secretLoader.LoadSecret(context.Background(), "BIBLE_API_URL")
		bibleAPIKey, _ = secretLoader.LoadSecret(context.Background(), "BIBLE_API_KEY")
		llmSystemPrompts, _ = secretLoader.LoadSecret(context.Background(), "LLM_SYSTEM_PROMPTS")
	} else {
		bibleAPIURL = os.Getenv("BIBLE_API_URL")
		bibleAPIKey = os.Getenv("BIBLE_API_KEY")
		llmSystemPrompts = os.Getenv("LLM_SYSTEM_PROMPTS")
	}

	if bibleAPIURL != "" && bibleAPIKey != "" {
		bibleAIClient = services.NewRealBibleAIClient(bibleAPIURL, bibleAPIKey, llmSystemPrompts)
		logger.Info("Using RealBibleAIClient", "url", bibleAPIURL)
	} else {
		logger.Info("BIBLE_API_URL or BIBLE_API_KEY not set, using MockBibleAIClient")
		bibleAIClient = services.NewMockBibleAIClient()
	}

	// Create a safe database wrapper that handles nil database connection
	var dbWrapper services.DBInterface
	if database.DB != nil {
		dbWrapper = database.DB
	} else {
		// Create a mock database that returns errors for all operations
		dbWrapper = &mockDatabase{connected: false}
	}

	noteService := services.NewNoteService(dbWrapper)

	bibleHandler := handlers.NewBibleHandler(bibleAIClient)
	chatHandler := handlers.NewChatHandler(bibleAIClient, noteService, dbWrapper)
	noteHandler := handlers.NewNoteHandler(dbWrapper, noteService)
	userHandler := handlers.NewUserHandler(dbWrapper)

	// Update handlers to use notification service
	connectionHandler := handlers.NewConnectionHandler(dbWrapper, notificationService)
	groupHandler := handlers.NewGroupHandler(dbWrapper, notificationService)
	groupShareHandler := handlers.NewGroupShareHandler(dbWrapper, notificationService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Reading Plans
	readingPlanService := services.NewReadingPlanService(dbWrapper)
	readingPlanHandler := handlers.NewReadingPlanHandler(readingPlanService)

	// Memory Verses
	memoryVerseService := services.NewMemoryVerseService(dbWrapper)
	memoryVerseHandler := handlers.NewMemoryVerseHandler(memoryVerseService)

	// Study Templates
	templateService := services.NewTemplateService(dbWrapper, bibleAIClient)
	templateHandler := handlers.NewTemplateHandler(templateService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestLogger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			// Check env var first for explicit allowed origins
			if allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); allowedOrigins != "" {
				for _, allowed := range strings.Split(allowedOrigins, ",") {
					if strings.TrimSpace(allowed) == origin {
						return true
					}
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if !dbConnected {
			w.WriteHeader(http.StatusServiceUnavailable)
			logger.Error("Health check failed: DB not connected")
			if _, err := w.Write([]byte("DATABASE_CONNECTION_FAILED")); err != nil {
				logger.Error("Failed to write health check response", "error", err)
			}
			return
		}

		if err := database.DB.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			logger.Error("Health check failed: DB ping failed", "error", err)
			if _, err := w.Write([]byte("DATABASE_PING_FAILED")); err != nil {
				logger.Error("Failed to write health check response", "error", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("Failed to write health check response", "error", err)
		}
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Group(func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware)
		}
		r.Post("/api/users/me", userHandler.CreateOrUpdateUser)
		r.Put("/api/users/me", userHandler.UpdateUser)
		r.Get("/api/users/me", userHandler.GetMe)

		r.Get("/api/notes", noteHandler.GetNotes)
		r.Post("/api/notes", noteHandler.CreateNote)
		r.Get("/api/notes/{id}", noteHandler.GetNote)
		r.Put("/api/notes/{id}", noteHandler.UpdateNote)
		r.Delete("/api/notes/{id}", noteHandler.DeleteNote)

		r.Get("/api/bible/passage", bibleHandler.GetBiblePassage)
		r.Post("/api/chat", chatHandler.ChatWithAI)
		r.Post("/api/ai/ask", chatHandler.AskAI)

		// Notifications
		r.Post("/api/notifications/register", notificationHandler.RegisterDevice)

		// Connections
		r.Get("/api/users/search", connectionHandler.SearchUsers)
		r.Post("/api/connections/request", connectionHandler.SendConnectionRequest)
		r.Get("/api/connections", connectionHandler.ListConnections)
		r.Put("/api/connections/{id}", connectionHandler.AcceptConnectionRequest)
		r.Delete("/api/connections/{id}", connectionHandler.DeleteConnectionRequest)

		// Groups
		r.Post("/api/groups", groupHandler.CreateGroup)
		r.Get("/api/groups", groupHandler.ListMyGroups)
		r.Get("/api/groups/search", groupHandler.SearchGroups)
		r.Post("/api/groups/{id}/join", groupHandler.JoinGroup)
		r.Delete("/api/groups/{id}/leave", groupHandler.LeaveGroup)
		r.Get("/api/groups/{id}/members", groupHandler.GetGroupMembers)
		r.Post("/api/groups/{id}/members", groupHandler.AddGroupMember)
		r.Delete("/api/groups/{id}/members/{userId}", groupHandler.RemoveGroupMember)

		// Group Shares
		r.Post("/api/groups/{id}/shares", groupShareHandler.ShareNoteToGroup)
		r.Get("/api/groups/{id}/shares", groupShareHandler.ListGroupShares)
		r.Get("/api/groups/{id}/shares/{shareId}", groupShareHandler.GetSharedNoteDetails)

		// Reading Plans
		r.Get("/api/reading-plans", readingPlanHandler.GetAllPlans)
		r.Get("/api/reading-plans/{id}", readingPlanHandler.GetPlan)
		r.Post("/api/reading-plans/{id}/subscribe", readingPlanHandler.Subscribe)
		r.Get("/api/my-reading-plans", readingPlanHandler.GetUserPlans)
		r.Post("/api/my-reading-plans/{id}/progress", readingPlanHandler.MarkDayComplete)
		r.Get("/api/my-reading-plans/{id}/progress", readingPlanHandler.GetPlanProgress)

		// Memory Verses (Refactored)
		r.Get("/api/verse-packs", memoryVerseHandler.GetPacks)
		r.Post("/api/verse-packs", memoryVerseHandler.CreatePack)
		r.Get("/api/verse-packs/{id}", memoryVerseHandler.GetPackDetails)
		r.Delete("/api/verse-packs/{id}", memoryVerseHandler.DeletePack)
		r.Post("/api/verse-packs/{id}/verses", memoryVerseHandler.CreateVerseInPack)
		r.Post("/api/verse-packs/{id}/clone", memoryVerseHandler.ClonePack)

		// Study Templates
		r.Post("/api/templates", templateHandler.CreateTemplate)
		r.Get("/api/templates", templateHandler.ListMyTemplates)
		r.Get("/api/templates/public", templateHandler.ListPublicTemplates)
		r.Get("/api/templates/{id}", templateHandler.GetTemplate)
		r.Put("/api/templates/{id}", templateHandler.UpdateTemplate)
		r.Delete("/api/templates/{id}", templateHandler.DeleteTemplate)
		r.Post("/api/templates/{id}/clone", templateHandler.CloneTemplate)
		r.Post("/api/templates/{id}/generate", templateHandler.Generate)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		logger.Info(fmt.Sprintf("Server listening on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exiting")
}
