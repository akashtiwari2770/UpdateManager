package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// ReadJSON reads and decodes JSON from request body
func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Limit request body size (10MB)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Check for extra data
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return err
	}

	return nil
}

// GetQueryParam gets a query parameter with default value
func GetQueryParam(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}

// GetIntQueryParam gets an integer query parameter with default value
func GetIntQueryParam(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
