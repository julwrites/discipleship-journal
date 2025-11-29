package main

import (
	"context"
	"fmt"
	"log"
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
)

func main() {
	// Load .env file if it exists (local dev)
	_ = godotenv.Load()

	if err := database.Connect(); err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	}
	defer database.Close()

	// Init Auth Middleware (assumes GOOGLE_APPLICATION_CREDENTIALS or similar set in prod)
	// For local dev with a specific key file, pass the content or path.
	// Here we use default context.
	authMiddleware, err := middleware.NewAuthMiddleware(context.Background())
	if err != nil {
		log.Printf("Warning: Firebase Auth init failed: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware.VerifyToken)
		}
		r.Post("/api/users/me", handlers.CreateOrUpdateUser)
		r.Put("/api/users/me", handlers.UpdateUser)

		r.Get("/api/notes", handlers.GetNotes)
		r.Post("/api/notes", handlers.CreateNote)
		r.Get("/api/notes/{id}", handlers.GetNote)
		r.Put("/api/notes/{id}", handlers.UpdateNote)

		r.Get("/api/bible/passage", handlers.GetBiblePassage)
		r.Post("/api/chat", handlers.ChatWithAI)
		r.Post("/api/ai/ask", handlers.AskAI)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("Server listening on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
