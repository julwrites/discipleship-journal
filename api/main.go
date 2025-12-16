package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	if err := database.Connect(); err != nil {
		logger.Error("Database connection failed", "error", err)
	}
	defer database.Close()

	// Init Firebase Service
	// Pass empty string for saKey to use default credentials (production)
	// or rely on GOOGLE_APPLICATION_CREDENTIALS
	firebaseService, err := services.NewFirebaseService(context.Background(), "", "discipleship-journal-pwa")
	if err != nil {
		logger.Error("Firebase init failed", "error", err)
	}

	var authMiddleware *middleware.AuthMiddleware
	var notificationService services.NotificationService

	if firebaseService != nil {
		authMiddleware = middleware.NewAuthMiddlewareFromClient(firebaseService.AuthClient)
		if firebaseService.MessagingClient != nil {
			notificationService = services.NewNotificationService(database.DB, firebaseService.MessagingClient)
		} else {
			logger.Warn("Firebase Messaging not initialized")
			notificationService = services.NewMockNotificationService() // Fallback to avoid nil pointer
		}
	} else {
		// Fallback for local dev without firebase creds?
		// We can't really auth without firebase.
		// Existing code allowed running but AuthMiddleware would fail.
		logger.Error("Firebase Service is nil")
	}

	// Init Bible AI Client
	var bibleAIClient services.BibleAIClient
	bibleAPIURL := os.Getenv("BIBLE_API_URL")

	if bibleAPIURL != "" {
		bibleAIClient = services.NewRealBibleAIClient(bibleAPIURL, os.Getenv("BIBLE_API_KEY"))
	} else {
		logger.Info("BIBLE_API_URL not set, using MockBibleAIClient")
		bibleAIClient = services.NewMockBibleAIClient()
	}

	noteService := services.NewNoteService(database.DB)

	bibleHandler := handlers.NewBibleHandler(bibleAIClient)
	chatHandler := handlers.NewChatHandler(bibleAIClient, noteService)
	noteHandler := handlers.NewNoteHandler(database.DB, noteService)

	// Update handlers to use notification service
	connectionHandler := handlers.NewConnectionHandler(database.DB, notificationService)
	groupHandler := handlers.NewGroupHandler(database.DB, notificationService)
	groupShareHandler := handlers.NewGroupShareHandler(database.DB, notificationService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestLogger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.DB.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			logger.Error("Health check failed: DB not connected", "error", err)
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
			r.Use(authMiddleware.VerifyToken)
		}
		r.Post("/api/users/me", handlers.CreateOrUpdateUser)
		r.Put("/api/users/me", handlers.UpdateUser)

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
