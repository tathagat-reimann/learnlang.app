package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{"message": "hello"}

	WriteJSON(rr, http.StatusOK, payload)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	var result map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["message"] != "hello" {
		t.Errorf("Expected message 'hello', got %s", result["message"])
	}
}

func TestWriteErrorCode(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteErrorCode(rr, http.StatusBadRequest, "INVALID", "Invalid input")

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if resp.Code != "INVALID" || resp.Error != "Invalid input" {
		t.Errorf("Unexpected error response: %+v", resp)
	}
}

func TestWriteOKData(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"result": "ok"}
	meta := map[string]string{"version": "1.0"}

	WriteOKData(rr, data, meta)

	var resp SuccessResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse success response: %v", err)
	}

	d := resp.Data.(map[string]interface{})
	if d["result"] != "ok" {
		t.Errorf("Expected result 'ok', got %v", d["result"])
	}
}

func TestGetRequestID(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "abc-123")

	id := GetRequestID(req)
	if id != "abc-123" {
		t.Errorf("Expected request ID 'abc-123', got %s", id)
	}
}

func TestWriteErrorWithRequest_Full(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "req-789")

	expected := ErrorResponse{
		Error:     "Input is not valid",
		Code:      "INVALID_INPUT",
		RequestID: "req-789",
	}

	WriteErrorWithRequest(rr, req, 400, "INVALID_INPUT", "Input is not valid")

	assertHTTPResponse(t, rr, 400, "application/json", expected)
}

// assertHTTPResponse checks status code, headers, and JSON body.
func assertHTTPResponse[T any](t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedContentType string, expectedBody T) {
	t.Helper()

	// Status code
	if rr.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, rr.Code)
	}

	// Content-Type header
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("Expected Content-Type %q, got %q", expectedContentType, ct)
	}

	// JSON body
	var actual T
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !reflect.DeepEqual(actual, expectedBody) {
		t.Errorf("Expected body %+v, got %+v", expectedBody, actual)
	}
}
