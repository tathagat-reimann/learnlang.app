package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"learnlang-backend/router"
	"learnlang-backend/store"
)

// tinyPNG returns a minimal PNG header to satisfy content sniffing (image/png)
func tinyPNG() []byte {
	return []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
}

func newMultipartVocabReq(t *testing.T, url, name, packID string) (*http.Request, error) {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	// fields
	if err := w.WriteField("name", name); err != nil {
		return nil, err
	}
	// Provide translation to satisfy required validation; use name as default
	if err := w.WriteField("translation", name); err != nil {
		return nil, err
	}
	if err := w.WriteField("pack_id", packID); err != nil {
		return nil, err
	}
	// file
	fw, err := w.CreateFormFile("image", "file.png")
	if err != nil {
		return nil, err
	}
	if _, err := fw.Write(tinyPNG()); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}

func TestCreateVocab_And_GetFlashcards(t *testing.T) {
	os.Setenv("UPLOAD_DIR", t.TempDir())
	if err := store.InitFromEnv(); err != nil {
		t.Fatalf("db init failed: %v", err)
	}
	store.Reset()
	h := router.NewRouter()

	// Create a pack first
	pbody, _ := json.Marshal(createPackReq{Name: "Kitchen", LangID: "1", UserID: "u1"})
	pw := httptest.NewRecorder()
	h.ServeHTTP(pw, httptest.NewRequest(http.MethodPost, "/api/packs", bytes.NewReader(pbody)))
	if pw.Code != http.StatusCreated {
		t.Fatalf("pack create failed: %d %s", pw.Code, pw.Body.String())
	}

	// Add vocabs
	// Find the created pack ID using composite key
	packID := store.GetPackIDByKey("u1:1:kitchen")
	if packID == "" {
		t.Fatalf("pack not created")
	}
	for _, nm := range []string{"knife", "utensil"} {
		w := httptest.NewRecorder()
		req, err := newMultipartVocabReq(t, "/api/vocabs", nm, packID)
		if err != nil {
			t.Fatalf("multipart req err: %v", err)
		}
		h.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("vocab create failed: %d %s", w.Code, w.Body.String())
		}
		// assert file exists on disk
		var created struct {
			Data struct {
				Image string `json:"image"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("invalid json response: %v", err)
		}
		if created.Data.Image == "" || !strings.HasPrefix(created.Data.Image, "/files/") {
			t.Fatalf("unexpected image url: %q", created.Data.Image)
		}
		fname := strings.TrimPrefix(created.Data.Image, "/files/")
		diskPath := filepath.Join(os.Getenv("UPLOAD_DIR"), fname)
		if st, err := os.Stat(diskPath); err != nil || st.Size() <= 0 {
			t.Fatalf("expected file to exist at %s (err=%v)", diskPath, err)
		}
	}

	// Get flashcards
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/flashcards?user_id=u1&lang_id=1&pack_ids="+packID+"&limit=10", nil)
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("flashcards failed: %d %s", rw.Code, rw.Body.String())
	}
	var resp struct {
		Data []map[string]any `json:"data"`
		Meta map[string]any   `json:"meta"`
	}
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(resp.Data))
	}
}

func TestCreateVocab_Duplicate(t *testing.T) {
	os.Setenv("UPLOAD_DIR", t.TempDir())
	if err := store.InitFromEnv(); err != nil {
		t.Fatalf("db init failed: %v", err)
	}
	store.Reset()
	h := router.NewRouter()

	// Create pack
	pbody, _ := json.Marshal(createPackReq{Name: "Sports", LangID: "2", UserID: "u1"})
	pw := httptest.NewRecorder()
	h.ServeHTTP(pw, httptest.NewRequest(http.MethodPost, "/api/packs", bytes.NewReader(pbody)))
	if pw.Code != http.StatusCreated {
		t.Fatalf("pack create failed: %d %s", pw.Code, pw.Body.String())
	}

	// Create vocab
	// lookup packID via composite key for Sports/de/u1
	sportID := store.GetPackIDByKey("u1:2:sports")
	if sportID == "" {
		t.Fatalf("sports pack not found")
	}
	w1 := httptest.NewRecorder()
	req1, err := newMultipartVocabReq(t, "/api/vocabs", "ball", sportID)
	if err != nil {
		t.Fatalf("multipart req err: %v", err)
	}
	h.ServeHTTP(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("vocab create failed: %d %s", w1.Code, w1.Body.String())
	}
	// assert file exists for first create
	var created struct {
		Data struct {
			Image string `json:"image"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w1.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	fname := strings.TrimPrefix(created.Data.Image, "/files/")
	diskPath := filepath.Join(os.Getenv("UPLOAD_DIR"), fname)
	if st, err := os.Stat(diskPath); err != nil || st.Size() <= 0 {
		t.Fatalf("expected file to exist at %s (err=%v)", diskPath, err)
	}

	// Duplicate
	w2 := httptest.NewRecorder()
	req2, err := newMultipartVocabReq(t, "/api/vocabs", "ball", sportID)
	if err != nil {
		t.Fatalf("multipart req err: %v", err)
	}
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w2.Code, w2.Body.String())
	}
}
