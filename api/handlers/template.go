package handlers

import (
	"encoding/json"
	"net/http"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"firebase.google.com/go/v4/auth"
)

type TemplateHandler struct {
	service services.TemplateService
}

func NewTemplateHandler(service services.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tmpl models.StudyTemplate
	if err := json.NewDecoder(r.Body).Decode(&tmpl); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	tmpl.CreatorID = userID

	created, err := h.service.CreateTemplate(r.Context(), &tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *TemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tmpl, err := h.service.GetTemplate(r.Context(), id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tmpl)
}

func (h *TemplateHandler) ListMyTemplates(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	list, err := h.service.ListTemplates(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": list})
}

func (h *TemplateHandler) ListPublicTemplates(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ListPublicTemplates(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": list})
}

func (h *TemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var tmpl models.StudyTemplate
	if err := json.NewDecoder(r.Body).Decode(&tmpl); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	tmpl.ID = id
	tmpl.CreatorID = userID

	if err := h.service.UpdateTemplate(r.Context(), &tmpl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTemplate(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *TemplateHandler) CloneTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cloned, err := h.service.CloneTemplate(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloned)
}

func (h *TemplateHandler) Generate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	content, err := h.service.GenerateContent(r.Context(), id, req.Inputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"content": content})
}

func (h *TemplateHandler) getUserID(r *http.Request) (uuid.UUID, error) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		firebaseUID = token.UID
	}
	return GetUserUUID(r.Context(), firebaseUID)
}
