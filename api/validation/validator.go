package validation

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Register function to get json tag name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns validation errors if any
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// DecodeAndValidate decodes a JSON request body into the provided struct and validates it.
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return false
	}

	if err := validate.Struct(v); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			http.Error(w, "Validation error", http.StatusInternalServerError)
			return false
		}

		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = "Validation failed on the '" + err.Tag() + "' tag"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
		return false
	}

	return true
}
