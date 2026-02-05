package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminHandler_SyncBibleVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBibleVersionService)
		handler := NewAdminHandler(mockService)

		mockService.On("SyncVersions", mock.Anything).Return(nil)

		req, _ := http.NewRequest("POST", "/admin/sync-bible-versions", nil)
		rr := httptest.NewRecorder()

		handler.SyncBibleVersions(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

        var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, true, response["success"])
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockBibleVersionService)
		handler := NewAdminHandler(mockService)

		mockService.On("SyncVersions", mock.Anything).Return(errors.New("sync failed"))

		req, _ := http.NewRequest("POST", "/admin/sync-bible-versions", nil)
		rr := httptest.NewRecorder()

		handler.SyncBibleVersions(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
