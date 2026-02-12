package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMemoryVerseHandler_CreatePack_InvalidBody(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("POST", "/api/verse-packs", strings.NewReader("invalid json"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreatePack(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_CreatePack_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	payload := `{"title": "New Pack"}`
	req := httptest.NewRequest("POST", "/api/verse-packs", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("CreatePack", mock.Anything, mock.Anything).Return(&models.VersePack{}, errors.New("service error"))

	handler.CreatePack(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_GetPackDetails_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Get("/api/verse-packs/{id}", handler.GetPackDetails)

	req := httptest.NewRequest("GET", "/api/verse-packs/invalid-uuid", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_GetPackDetails_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Get("/api/verse-packs/{id}", handler.GetPackDetails)

	req := httptest.NewRequest("GET", "/api/verse-packs/"+packID.String(), nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("GetPack", mock.Anything, packID, userID).Return(nil, errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_CreateVerseInPack_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/verses", handler.CreateVerseInPack)

	req := httptest.NewRequest("POST", "/api/verse-packs/invalid-uuid/verses", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_CreateVerseInPack_PackNotFound(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/verses", handler.CreateVerseInPack)

	req := httptest.NewRequest("POST", "/api/verse-packs/"+packID.String()+"/verses", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("GetPack", mock.Anything, packID, userID).Return(nil, errors.New("not found"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMemoryVerseHandler_CreateVerseInPack_Forbidden(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	otherUserID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/verses", handler.CreateVerseInPack)

	req := httptest.NewRequest("POST", "/api/verse-packs/"+packID.String()+"/verses", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Pack owned by someone else
	pack := &models.VersePack{ID: packID, UserID: &otherUserID}
	mockService.On("GetPack", mock.Anything, packID, userID).Return(pack, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMemoryVerseHandler_CreateVerseInPack_InvalidBody(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/verses", handler.CreateVerseInPack)

	req := httptest.NewRequest("POST", "/api/verse-packs/"+packID.String()+"/verses", strings.NewReader("invalid json"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	pack := &models.VersePack{ID: packID, UserID: &userID}
	mockService.On("GetPack", mock.Anything, packID, userID).Return(pack, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_CreateVerseInPack_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/verses", handler.CreateVerseInPack)

	payload := `{"reference": "John 3:16"}`
	req := httptest.NewRequest("POST", "/api/verse-packs/"+packID.String()+"/verses", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	pack := &models.VersePack{ID: packID, UserID: &userID}
	mockService.On("GetPack", mock.Anything, packID, userID).Return(pack, nil)
	mockService.On("CreateVerse", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_ClonePack_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/clone", handler.ClonePack)

	req := httptest.NewRequest("POST", "/api/verse-packs/invalid-uuid/clone", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_ClonePack_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/clone", handler.ClonePack)

	req := httptest.NewRequest("POST", "/api/verse-packs/"+packID.String()+"/clone", strings.NewReader("{}"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("ClonePack", mock.Anything, packID, userID, "", true).Return(nil, errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_DeletePack_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/verse-packs/{id}", handler.DeletePack)

	req := httptest.NewRequest("DELETE", "/api/verse-packs/invalid-uuid", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_DeletePack_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/verse-packs/{id}", handler.DeletePack)

	req := httptest.NewRequest("DELETE", "/api/verse-packs/"+packID.String(), nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("DeletePack", mock.Anything, packID, userID).Return(errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_UpdateVerse_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}", handler.UpdateVerse)

	req := httptest.NewRequest("PUT", "/api/memory-verses/invalid-uuid", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_UpdateVerse_InvalidBody(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}", handler.UpdateVerse)

	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String(), strings.NewReader("invalid json"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_UpdateVerse_NotFound(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}", handler.UpdateVerse)

	payload := `{"reference": "John 3:17"}`
	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String(), strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("UpdateVerse", mock.Anything, mock.Anything, userID).Return(models.ErrNotFound)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMemoryVerseHandler_UpdateVerse_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}", handler.UpdateVerse)

	payload := `{"reference": "John 3:17"}`
	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String(), strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("UpdateVerse", mock.Anything, mock.Anything, userID).Return(errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_DeleteVerse_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}", handler.DeleteVerse)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/invalid-uuid", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_DeleteVerse_NotFound(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}", handler.DeleteVerse)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/"+verseID.String(), nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("DeleteVerse", mock.Anything, verseID, userID).Return(models.ErrNotFound)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMemoryVerseHandler_DeleteVerse_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}", handler.DeleteVerse)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/"+verseID.String(), nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("DeleteVerse", mock.Anything, verseID, userID).Return(errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_SetVersePreference_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}/preference", handler.SetVersePreference)

	req := httptest.NewRequest("PUT", "/api/memory-verses/invalid-uuid/preference", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_SetVersePreference_InvalidBody(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}/preference", handler.SetVersePreference)

	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String()+"/preference", strings.NewReader("invalid json"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_SetVersePreference_MissingVersion(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}/preference", handler.SetVersePreference)

	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String()+"/preference", strings.NewReader("{}"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_SetVersePreference_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}/preference", handler.SetVersePreference)

	payload := `{"version": "NIV"}`
	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String()+"/preference", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("SetVersePreference", mock.Anything, userID, verseID, "NIV").Return(errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_RemoveVersePreference_InvalidID(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}/preference", handler.RemoveVersePreference)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/invalid-uuid/preference", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMemoryVerseHandler_RemoveVersePreference_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}/preference", handler.RemoveVersePreference)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/"+verseID.String()+"/preference", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("RemoveVersePreference", mock.Anything, userID, verseID).Return(errors.New("service error"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_SearchVerses_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("GET", "/api/memory-verses?q=test", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("SearchVerses", mock.Anything, userID, "test").Return(nil, errors.New("service error"))

	handler.SearchVerses(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_GetPacks_ServiceError(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("GET", "/api/verse-packs", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("GetPacks", mock.Anything, userID, "user").Return(nil, errors.New("service error"))

	handler.GetPacks(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMemoryVerseHandler_GetPackDetails_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	r := chi.NewRouter()
	r.Get("/api/verse-packs/{id}", handler.GetPackDetails)

	req := httptest.NewRequest("GET", "/api/verse-packs/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMemoryVerseHandler_DeletePack_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/api/verse-packs/{id}", handler.DeletePack)

	req := httptest.NewRequest("DELETE", "/api/verse-packs/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMemoryVerseHandler_SearchVerses_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	req := httptest.NewRequest("GET", "/api/memory-verses?q=test", nil)
	w := httptest.NewRecorder()

	handler.SearchVerses(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMemoryVerseHandler_UpdateVerse_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}", handler.UpdateVerse)

	req := httptest.NewRequest("PUT", "/api/memory-verses/"+uuid.NewString(), strings.NewReader("{}"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMemoryVerseHandler_DeleteVerse_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}", handler.DeleteVerse)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMemoryVerseHandler_SetVersePreference_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}/preference", handler.SetVersePreference)

	req := httptest.NewRequest("PUT", "/api/memory-verses/"+uuid.NewString()+"/preference", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMemoryVerseHandler_RemoveVersePreference_Unauthorized(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}/preference", handler.RemoveVersePreference)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/"+uuid.NewString()+"/preference", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
