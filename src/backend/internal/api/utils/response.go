package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a JSON response
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MetaInfo represents pagination or metadata
type MetaInfo struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a success JSON response
func WriteSuccess(w http.ResponseWriter, status int, data interface{}) error {
	response := JSONResponse{
		Success: true,
		Data:    data,
	}
	return WriteJSON(w, status, response)
}

// WriteError writes an error JSON response
func WriteError(w http.ResponseWriter, status int, code, message string) error {
	response := JSONResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
	return WriteJSON(w, status, response)
}

// WritePaginated writes a paginated JSON response
func WritePaginated(w http.ResponseWriter, status int, data interface{}, page, limit int, total int64) error {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := JSONResponse{
		Success: true,
		Data:    data,
		Meta: &MetaInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
	return WriteJSON(w, status, response)
}
