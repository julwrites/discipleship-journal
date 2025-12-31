package contract

import (
	"context"
	"net/http"
	"testing"

	"discipleship_journal_api/handlers"
	"discipleship_journal_api/services"
	"discipleship_journal_api/tests/integration"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupContractTest spins up the integration DB, sets up the router with real services,
// and returns the router, the pool, and a cleanup function.
func SetupContractTest(t *testing.T) (*chi.Mux, *pgxpool.Pool, func()) {
	pool, cleanup := integration.SetupIntegrationDB(t)

	// Initialize Services
	noteService := services.NewNoteService(pool)

	// Initialize Handlers
	noteHandler := handlers.NewNoteHandler(pool, noteService)

	// Setup Router
	r := chi.NewRouter()

	testAuthMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Inject a test user into context.
			// Using a valid UUID to ensure DB compatibility.
			ctx := context.WithValue(r.Context(), handlers.TestUserKey, "00000000-0000-0000-0000-000000000001")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r.Use(testAuthMiddleware)

	r.Route("/api/notes", func(r chi.Router) {
		r.Get("/", noteHandler.GetNotes)
		r.Post("/", noteHandler.CreateNote)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", noteHandler.GetNote)
			r.Put("/", noteHandler.UpdateNote)
			r.Delete("/", noteHandler.DeleteNote)
		})
	})

	return r, pool, cleanup
}
