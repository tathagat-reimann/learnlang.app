package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// ErrorResponse standard JSON error body.
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// WriteJSON writes a JSON response with status code and content-type header.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteOK writes a 200 OK JSON response.
func WriteOK(w http.ResponseWriter, v any) {
	WriteJSON(w, http.StatusOK, v)
}

// SuccessResponse is a standard envelope for successful responses.
type SuccessResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

// WriteData writes a JSON success envelope with the given status.
func WriteData(w http.ResponseWriter, status int, data any, meta any) {
	WriteJSON(w, status, SuccessResponse{Data: data, Meta: meta})
}

// WriteOKData writes a 200 OK success envelope.
func WriteOKData(w http.ResponseWriter, data any, meta any) {
	WriteData(w, http.StatusOK, data, meta)
}

// WriteCreatedData writes a 201 Created success envelope.
func WriteCreatedData(w http.ResponseWriter, data any, meta any) {
	WriteData(w, http.StatusCreated, data, meta)
}

// GetRequestID retrieves a request ID if available.
func GetRequestID(r *http.Request) string {
	if r == nil {
		return ""
	}
	if id := middleware.GetReqID(r.Context()); id != "" {
		return id
	}
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return ""
}

// WriteError writes a standardized JSON error body with message only.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}

// WriteErrorCode writes a standardized JSON error body with an error code.
func WriteErrorCode(w http.ResponseWriter, status int, code, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg, Code: code})
}

// WriteErrorWithRequest writes a standardized JSON error with code and request ID.
func WriteErrorWithRequest(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	log.Printf("Error [%s]: %s (RequestID: %s)", code, msg, GetRequestID(r))
	WriteJSON(w, status, ErrorResponse{Error: msg, Code: code, RequestID: GetRequestID(r)})
}
