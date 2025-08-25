package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"learnlang-backend/router"
	"learnlang-backend/store"
)

type createPackReq struct {
	Name   string `json:"name"`
	LangID string `json:"lang_id"`
	UserID string `json:"user_id"`
}

func setup(t *testing.T) http.Handler {
	t.Helper()
	// Ensure uploads dir env is set for router/static and handlers
	os.Setenv("UPLOAD_DIR", t.TempDir())
	if err := store.InitFromEnv(); err != nil {
		t.Fatalf("db init failed: %v", err)
	}
	store.Reset()
	return router.NewRouter()
}

func TestCreatePack_Success(t *testing.T) {
	h := setup(t)

	body, _ := json.Marshal(createPackReq{Name: "Basics", LangID: "1", UserID: "u1"})
	req := httptest.NewRequest(http.MethodPost, "/api/packs", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}

	// ensure envelope has data
	var resp struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Data["id"] == nil || resp.Data["name"] != "Basics" {
		t.Fatalf("unexpected data: %#v", resp.Data)
	}
}

func TestCreatePack_Duplicate(t *testing.T) {
	h := setup(t)

	body, _ := json.Marshal(createPackReq{Name: "Basics", LangID: "1", UserID: "u1"})
	req1 := httptest.NewRequest(http.MethodPost, "/api/packs", bytes.NewReader(body))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("setup create failed: %d %s", w1.Code, w1.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "/api/packs", bytes.NewReader(body))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w2.Code)
	}
}

func TestGetLanguages_OK(t *testing.T) {
	h := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/api/languages", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp.Data) == 0 {
		t.Fatalf("expected languages, got none")
	}
}
