package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `json:"name" validate:"required"`
}

func TestDecodeAndValidate_Coverage(t *testing.T) {
	t.Run("JSON_Decode_Error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()
		var s TestStruct

		ok := DecodeAndValidate(w, req, &s)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("Validation_Error_Required", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader("{}")) // Name missing
		w := httptest.NewRecorder()
		var s TestStruct

		ok := DecodeAndValidate(w, req, &s)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		errors, _ := resp["errors"].(map[string]interface{})
		assert.Contains(t, errors, "name")
	})

	t.Run("Invalid_Validation_Type", func(t *testing.T) {
		// Passing a non-pointer to Decode (which fails json decode first actually)
		// Wait, Decode requires a pointer. If we pass non-pointer to DecodeAndValidate:
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))
		w := httptest.NewRecorder()
		var s TestStruct

		// Decode will return error: "json: Unmarshal(non-pointer ...)"
		// So it hits the first error block.
		ok := DecodeAndValidate(w, req, s)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid_Validation_Structure", func(t *testing.T) {
		// To trigger validator.InvalidValidationError, we need to pass something that json.Decode accepts
		// but validator rejects (like a nil pointer, or maybe a pointer to nil?).
		// But DecodeAndValidate calls json.Decode first.
		// If we pass a valid pointer to struct, json.Decode works.
		// validator.ValidateStruct(v) expects a struct or pointer to struct.
		// If we pass a map, validator might complain? No, it validates map?
		// Actually the project uses a wrapper `validation.ValidateStruct`.
		// Let's assume standard behavior.
		// If we pass `(*TestStruct)(nil)`, json.Decode might fail or work.
		// Let's try passing a pointer to an interface that holds nil?
		// Generally hard to hit InvalidValidationError if json.Decode succeeded, unless it's a type validator doesn't handle but json does.
		// e.g. `int`.
		req := httptest.NewRequest("POST", "/", strings.NewReader("123"))
		w := httptest.NewRecorder()
		var i int

		// json decodes 123 into i.
		// validation.ValidateStruct(&i) -> validator checks int.
		// validator.Struct(v) expects a struct.
		// It returns InvalidValidationError if not struct.

		ok := DecodeAndValidate(w, req, &i)
		assert.False(t, ok)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
