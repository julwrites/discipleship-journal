package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestDecodeAndValidate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		reqBody := TestRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		var v TestRequest
		valid := DecodeAndValidate(w, req, &v)

		assert.True(t, valid)
		assert.Equal(t, reqBody.Name, v.Name)
		assert.Equal(t, reqBody.Email, v.Email)
		assert.Equal(t, http.StatusOK, w.Code) // Default recorder code is 200 unless written to
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		var v TestRequest
		valid := DecodeAndValidate(w, req, &v)

		assert.False(t, valid)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("Validation Error", func(t *testing.T) {
		reqBody := TestRequest{
			Name:  "", // Missing name
			Email: "invalid-email",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		var v TestRequest
		valid := DecodeAndValidate(w, req, &v)

		assert.False(t, valid)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Validation failed")
		assert.Contains(t, w.Body.String(), "name")  // Lowercase due to RegisterTagNameFunc
		assert.Contains(t, w.Body.String(), "email") // Lowercase due to RegisterTagNameFunc
	})

	t.Run("Invalid Validation Error (Non-Struct)", func(t *testing.T) {
		// Valid JSON for int
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("123"))
		w := httptest.NewRecorder()

		var i int
		// DecodeAndValidate calls ValidateStruct(&i), which calls validator.Struct(*int).
		// Validator requires a struct, so it returns InvalidValidationError.
		valid := DecodeAndValidate(w, req, &i)

		assert.False(t, valid)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Validation error")
	})
}
