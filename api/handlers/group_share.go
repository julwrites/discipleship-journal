package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"discipleship_journal_api/services"
	chi "github.com/go-chi/chi/v5"
	pgx "github.com/jackc/pgx/v5"
)

type GroupShareHandler struct {
	db                  DBInterface
	notificationService services.NotificationService
}

func NewGroupShareHandler(db DBInterface, notificationService services.NotificationService) *GroupShareHandler {
	return &GroupShareHandler{db: db, notificationService: notificationService}
}

type ShareItemRequest struct {
	NoteID      *string `json:"note_id"`
	VersePackID *string `json:"verse_pack_id"`
	Comment     string  `json:"comment" validate:"max=500"`
}

type SharedItemResponse struct {
	ID          string      `json:"id"`
	GroupID     string      `json:"group_id"`
	NoteID      *string     `json:"note_id,omitempty"`
	VersePackID *string     `json:"verse_pack_id,omitempty"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle,omitempty"` // For Identifier or similar
	Content     interface{} `json:"content,omitempty"`  // For Notes
	Type        string      `json:"type"`               // "note" or "verse_pack"
	SharedBy    string      `json:"shared_by"`          // User display name
	SharedAt    string      `json:"shared_at"`
	Comment     string      `json:"comment"`
}

// ShareItemToGroup shares a note OR verse pack to a group
func (h *GroupShareHandler) ShareItemToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	var req ShareItemRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	if req.NoteID == nil && req.VersePackID == nil {
		http.Error(w, "Either note_id or verse_pack_id is required", http.StatusBadRequest)
		return
	}

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// 1. Verify membership
	var isMember bool
	err = h.db.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&isMember)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var resourceType string
	var resourceID string
	var resourceTitle string

	if req.NoteID != nil {
		resourceType = "note"
		resourceID = *req.NoteID
		// Verify ownership
		var ownerID, title string
		err = h.db.QueryRow(r.Context(),
			"SELECT user_id, title FROM notes WHERE id = $1 AND deleted_at IS NULL",
			resourceID).Scan(&ownerID, &title)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Note not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}
		if ownerID != userUUID.String() {
			http.Error(w, "You can only share your own notes", http.StatusForbidden)
			return
		}
		resourceTitle = title

		// Insert Share
		_, err = h.db.Exec(r.Context(),
			`INSERT INTO group_shares (group_id, note_id, shared_by, comment)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (group_id, note_id) WHERE note_id IS NOT NULL
			 DO UPDATE SET shared_at = NOW(), comment = $4`,
			groupID, resourceID, userUUID, req.Comment)
	} else if req.VersePackID != nil {
		resourceType = "verse_pack"
		resourceID = *req.VersePackID
		// Verify ownership OR public system pack
		var ownerID *string
		var title string
		var isPublic bool
		err = h.db.QueryRow(r.Context(),
			"SELECT user_id, title, is_public FROM verse_packs WHERE id = $1",
			resourceID).Scan(&ownerID, &title, &isPublic)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Verse pack not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Can share if owner OR is public
		canShare := isPublic || (ownerID != nil && *ownerID == userUUID.String())
		if !canShare {
			http.Error(w, "You cannot share this pack", http.StatusForbidden)
			return
		}
		resourceTitle = title

		// Insert Share
		// Note: The unique constraint on group_shares might need adjusting if we want to rely on ON CONFLICT.
		// Since we have separate columns and check constraints, we rely on a partial index if we want unique shares per type.
		// Assuming we don't have a unique constraint on verse_pack_id yet in the migration, we just insert.
		// Actually, standard practice for shares is unique per item.
		// For simplicity, I'll just INSERT. If the user shares again, it adds a new entry or we can unique it.
		// Let's assume we want to upsert if possible.
		// But Postgres unique constraint on (group_id, verse_pack_id) might not exist.
		// I will just perform an INSERT and ignore dupes or handle it.
		// Better: Check if already shared.
		var existingShareID string
		err = h.db.QueryRow(r.Context(),
			"SELECT id FROM group_shares WHERE group_id = $1 AND verse_pack_id = $2",
			groupID, resourceID).Scan(&existingShareID)
		if err == nil {
			// Update comment/time
			_, err = h.db.Exec(r.Context(), "UPDATE group_shares SET shared_at = NOW(), comment = $3 WHERE id = $1 AND group_id = $2", existingShareID, groupID, req.Comment)
		} else {
			_, err = h.db.Exec(r.Context(),
				"INSERT INTO group_shares (group_id, verse_pack_id, shared_by, comment) VALUES ($1, $2, $3, $4)",
				groupID, resourceID, userUUID, req.Comment)
		}
	}

	if err != nil {
		slog.Error("Failed to share item", "error", err)
		http.Error(w, "Failed to share item", http.StatusInternalServerError)
		return
	}

	// Notifications
	go h.sendNotifications(groupID, userUUID.String(), resourceTitle, resourceType, resourceID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *GroupShareHandler) sendNotifications(groupID, sharerID, resourceTitle, resourceType, resourceID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var groupName string
	if err := h.db.QueryRow(ctx, "SELECT name FROM groups WHERE id = $1", groupID).Scan(&groupName); err != nil {
		groupName = "Group"
	}

	var sharerName string
	if err := h.db.QueryRow(ctx, "SELECT COALESCE(username, email) FROM users WHERE id = $1", sharerID).Scan(&sharerName); err != nil {
		sharerName = "Someone"
	}

	rows, err := h.db.Query(ctx, "SELECT user_id FROM group_members WHERE group_id = $1 AND user_id != $2", groupID, sharerID)
	if err != nil {
		return
	}
	defer rows.Close()

	var memberIDs []string
	for rows.Next() {
		var memberID string
		if err := rows.Scan(&memberID); err == nil {
			memberIDs = append(memberIDs, memberID)
		}
	}

	if len(memberIDs) == 0 {
		return
	}

	msgBody := sharerName + " shared \"" + resourceTitle + "\" in " + groupName

	err = h.notificationService.SendMulticastNotification(ctx, memberIDs, "New Shared Item", msgBody, map[string]string{
		"type":        resourceType + "_share",
		"group_id":    groupID,
		"resource_id": resourceID,
	})
	if err != nil {
		slog.Error("Failed to send multicast notifications", "group_id", groupID, "error", err)
	}
}

// ListGroupShares lists items shared with the group
func (h *GroupShareHandler) ListGroupShares(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify membership
	var isMember bool
	err = h.db.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&isMember)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	query := `
		SELECT gs.id, gs.group_id, gs.note_id, gs.verse_pack_id,
		       COALESCE(n.title, vp.title) as title,
		       vp.identifier as subtitle,
		       COALESCE(u.username, u.email) as display_name, gs.shared_at, gs.comment,
		       CASE WHEN gs.note_id IS NOT NULL THEN 'note' ELSE 'verse_pack' END as type
		FROM group_shares gs
		LEFT JOIN notes n ON gs.note_id = n.id AND n.deleted_at IS NULL
		LEFT JOIN verse_packs vp ON gs.verse_pack_id = vp.id
		JOIN users u ON gs.shared_by = u.id
		WHERE gs.group_id = $1
		  AND (gs.note_id IS NULL OR n.id IS NOT NULL) -- Filter out deleted notes
		ORDER BY gs.shared_at DESC
	`

	rows, err := h.db.Query(r.Context(), query, groupID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shares []SharedItemResponse
	for rows.Next() {
		var s SharedItemResponse
		var sharedAt time.Time
		var subtitle *string
		var comment *string
		if err := rows.Scan(&s.ID, &s.GroupID, &s.NoteID, &s.VersePackID, &s.Title, &subtitle, &s.SharedBy, &sharedAt, &comment, &s.Type); err != nil {
			continue
		}
		if subtitle != nil {
			s.Subtitle = *subtitle
		}
		if comment != nil {
			s.Comment = *comment
		}
		s.SharedAt = sharedAt.Format(time.RFC3339)
		shares = append(shares, s)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shares); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetSharedItemDetails fetches details of a shared note or pack
func (h *GroupShareHandler) GetSharedItemDetails(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	shareID := chi.URLParam(r, "shareId")

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify membership
	var isMember bool
	err = h.db.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&isMember)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var s SharedItemResponse
	var sharedAt time.Time
	var subtitle *string
	var content interface{}
	var comment *string

	// Query to fetch generic details
	err = h.db.QueryRow(r.Context(),
		`SELECT gs.id, gs.group_id, gs.note_id, gs.verse_pack_id,
		        COALESCE(n.title, vp.title), vp.identifier, n.content,
		        COALESCE(u.username, u.email) as display_name, gs.shared_at, gs.comment,
		        CASE WHEN gs.note_id IS NOT NULL THEN 'note' ELSE 'verse_pack' END
		 FROM group_shares gs
		 LEFT JOIN notes n ON gs.note_id = n.id
		 LEFT JOIN verse_packs vp ON gs.verse_pack_id = vp.id
		 JOIN users u ON gs.shared_by = u.id
		 WHERE gs.id = $1 AND gs.group_id = $2`, shareID, groupID).Scan(
		&s.ID, &s.GroupID, &s.NoteID, &s.VersePackID, &s.Title, &subtitle, &content, &s.SharedBy, &sharedAt, &comment, &s.Type)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Shared item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if subtitle != nil {
		s.Subtitle = *subtitle
	}
	if comment != nil {
		s.Comment = *comment
	}
	s.Content = content
	s.SharedAt = sharedAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
