package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"discipleship_journal_api/validation"

	"github.com/go-playground/validator/v10"
)

// DecodeAndValidate decodes the request body into v and validates it.
// It writes an error response if decoding or validation fails.
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return false
	}

	if err := validation.ValidateStruct(v); err != nil {
		// Check if it's a validation error
		if _, ok := err.(*validator.InvalidValidationError); ok {
			slog.Error("Invalid validation error", "error", err)
			http.Error(w, "Validation error", http.StatusInternalServerError)
			return false
		}

		errors := make(map[string]string)
		// Use type assertion to get validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				errors[fieldErr.Field()] = "Validation failed on the '" + fieldErr.Tag() + "' tag"
			}
		} else {
			// Fallback for other errors
			errors["general"] = err.Error()
		}

		slog.Error("Validation failed", "errors", errors)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors}); err != nil {
			slog.Error("Failed to encode error response", "error", err)
		}
		return false
	}

	return true
}
