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
	// Mock Notification Service for tests
	notificationService := services.NewMockNotificationService()
	readingPlanService := services.NewReadingPlanService(pool)

	// Initialize Handlers
	noteHandler := handlers.NewNoteHandler(pool, noteService)
	groupService := services.NewGroupService(pool, notificationService)
	groupHandler := handlers.NewGroupHandler(groupService)
	groupShareHandler := handlers.NewGroupShareHandler(pool, notificationService)
	connectionHandler := handlers.NewConnectionHandler(pool, notificationService)
	readingPlanHandler := handlers.NewReadingPlanHandler(readingPlanService)

	// Setup Router
	r := chi.NewRouter()

	testAuthMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for override header
			userID := r.Header.Get("X-Test-User-ID")
			if userID == "" {
				userID = "00000000-0000-0000-0000-000000000001"
			}

			// Inject a test user into context.
			ctx := context.WithValue(r.Context(), handlers.TestUserKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r.Use(testAuthMiddleware)

	// Notes
	r.Route("/api/notes", func(r chi.Router) {
		r.Get("/", noteHandler.GetNotes)
		r.Post("/", noteHandler.CreateNote)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", noteHandler.GetNote)
			r.Put("/", noteHandler.UpdateNote)
			r.Delete("/", noteHandler.DeleteNote)
		})
	})

	// Groups
	r.Post("/api/groups", groupHandler.CreateGroup)
	r.Get("/api/groups", groupHandler.ListMyGroups)
	r.Get("/api/groups/search", groupHandler.SearchGroups)
	r.Post("/api/groups/direct", groupHandler.GetOrCreateDirectGroup) // Matches api.ts
	r.Route("/api/groups/{id}", func(r chi.Router) {
		r.Post("/join", groupHandler.JoinGroup)
		r.Delete("/leave", groupHandler.LeaveGroup)
		r.Get("/members", groupHandler.GetGroupMembers)
		r.Post("/members", groupHandler.AddGroupMember)
		r.Delete("/members/{userId}", groupHandler.RemoveGroupMember)
	})

	// Group Shares
	r.Post("/api/groups/{id}/shares", groupShareHandler.ShareItemToGroup)
	r.Get("/api/groups/{id}/shares", groupShareHandler.ListGroupShares)
	r.Get("/api/groups/{id}/shares/{shareId}", groupShareHandler.GetSharedItemDetails)


	// Connections
	r.Get("/api/users/search", connectionHandler.SearchUsers)
	r.Post("/api/connections/request", connectionHandler.SendConnectionRequest)
	r.Get("/api/connections", connectionHandler.ListConnections)
	r.Route("/api/connections/{id}", func(r chi.Router) {
		r.Put("/", connectionHandler.AcceptConnectionRequest)
		r.Delete("/", connectionHandler.DeleteConnectionRequest)
	})

	// Reading Plans
	r.Get("/api/reading-plans", readingPlanHandler.GetAllPlans)
	r.Get("/api/reading-plans/{id}", readingPlanHandler.GetPlan)
	r.Post("/api/reading-plans/{id}/subscribe", readingPlanHandler.Subscribe)

	r.Get("/api/my-reading-plans", readingPlanHandler.GetUserPlans)
	r.Route("/api/my-reading-plans/{id}/progress", func(r chi.Router) {
		r.Get("/", readingPlanHandler.GetPlanProgress)
		r.Post("/", readingPlanHandler.MarkDayComplete)
		r.Delete("/{day_number}", readingPlanHandler.UnmarkDayComplete)
	})

	return r, pool, cleanup
}
