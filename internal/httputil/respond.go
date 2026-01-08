package httputil

import (
	stdjson "encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// JSON writes JSON response with code
func JSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = stdjson.NewEncoder(w).Encode(v)
}

// JSONError writes an error message as JSON
func JSONError(w http.ResponseWriter, code int, message string) {
	JSON(w, code, map[string]string{"error": message})
}

// FormatValidationErrors converts validator.ValidationErrors into map[field]message
func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	m := make(map[string]string)
	for _, e := range errs {
		// Field returns the struct field name; use e.Field() as key
		m[e.Field()] = e.Error()
	}
	return m
}
