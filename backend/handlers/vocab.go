package handlers

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"learnlang-backend/models"
	"learnlang-backend/store"
	"learnlang-backend/utils"

	"github.com/google/uuid"
)

// CreateVocabRequestDTO represents request body to create a vocab.
// For multipart form uploads, we read from form fields instead of JSON body.
// This DTO remains for compatibility if you later accept JSON+URL uploads.
type CreateVocabRequestDTO struct {
	Image       string `json:"image"` // optional if using multipart file upload
	Name        string `json:"name"`
	Translation string `json:"translation"`
	PackID      string `json:"pack_id"`
}

// CreateVocabHandler creates a new vocab entry under a pack.
func CreateVocabHandler(w http.ResponseWriter, r *http.Request) {
	// Support multipart form for file uploads
	var (
		name        string
		translation string
		packID      string
		imgURL      string
	)
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		// Only multipart is supported now
		utils.WriteErrorWithRequest(w, r, http.StatusUnsupportedMediaType, utils.CodeInvalidJSON, "multipart/form-data required")
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidJSON, "invalid multipart form")
		return
	}
	name = strings.TrimSpace(r.FormValue("name"))
	translation = strings.TrimSpace(r.FormValue("translation"))
	packID = strings.TrimSpace(r.FormValue("pack_id"))

	file, header, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeMissingFields, "missing required field(s): image")
			return
		}
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidJSON, "invalid uploaded file")
		return
	}
	defer file.Close()

	// Enforce max size (10MB) and type (image/*) via sniffing
	const maxSize = 10 << 20 // 10MB
	// Quick size check from header if available
	if header.Size > 0 && header.Size > maxSize {
		utils.WriteErrorWithRequest(w, r, http.StatusRequestEntityTooLarge, utils.CodeFileTooLarge, "file too large; max 10MB")
		return
	}
	limited := http.MaxBytesReader(w, file, maxSize)
	// Read first 512 bytes to detect content type
	head := make([]byte, 512)
	n, err := io.ReadFull(limited, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) || errors.Is(err, io.EOF) {
			utils.WriteErrorWithRequest(w, r, http.StatusRequestEntityTooLarge, utils.CodeFileTooLarge, "file too large; max 10MB")
			return
		}
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidFileType, "could not read file header")
		return
	}
	contentType := http.DetectContentType(head[:n])
	if !strings.HasPrefix(contentType, "image/") {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidFileType, fmt.Sprintf("unsupported file type: %s", contentType))
		return
	}

	// Save file using utils.UploadImage
	url, err := utils.UploadImage(name, header.Filename, contentType, head, n, limited)
	if err != nil {
		// Map a few known errors to 4xx
		if strings.Contains(err.Error(), "unknown file type") {
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidFileType, "file type could not be determined")
			return
		}
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			utils.WriteErrorWithRequest(w, r, http.StatusRequestEntityTooLarge, utils.CodeFileTooLarge, "file too large; max 10MB")
			return
		}
		utils.WriteErrorWithRequest(w, r, http.StatusInternalServerError, utils.CodeInternal, "failed to save file")
		return
	}
	imgURL = url

	// Validate
	missing := make([]string, 0, 3)
	if packID == "" {
		missing = append(missing, "pack_id")
	}
	if imgURL == "" {
		missing = append(missing, "image")
	}
	if name == "" {
		missing = append(missing, "name")
	}
	if translation == "" {
		missing = append(missing, "translation")
	}
	if len(missing) > 0 {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeMissingFields, fmt.Sprintf("missing required field(s): %s", strings.Join(missing, ", ")))
		return
	}
	// pack must exist
	if _, ok := store.GetPackByID(packID); !ok {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidPack, fmt.Sprintf("unknown pack id: %q", packID))
		return
	}
	// uniqueness per pack/name
	vocabKey := utils.MakeVocabKeyByPackID(packID, name)
	if store.VocabExistsByKey(vocabKey) {
		utils.WriteErrorWithRequest(w, r, http.StatusConflict, utils.CodeDuplicateVocab, fmt.Sprintf("vocab %q already exists in this pack", name))
		return
	}

	v := models.Vocab{
		ID:          uuid.New().String(),
		Image:       imgURL,
		Name:        name,
		Translation: translation,
		PackID:      packID,
	}
	store.CreateVocab(v, vocabKey)
	utils.WriteCreatedData(w, v, nil)
}

// Flashcard represents a simplified view for the game (hide translation by default on UI).
type Flashcard struct {
	ID       string `json:"id"`
	Image    string `json:"image"`
	Name     string `json:"name"`
	PackName string `json:"pack_name"`
}

// GetFlashcardsHandler returns randomized flashcards for a user and language
// Optional query: packs=pack1,pack2 and limit=n (defaults to all and 20 max)
func GetFlashcardsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := strings.TrimSpace(q.Get("user_id"))
	lang := strings.TrimSpace(q.Get("lang_id"))
	if userID == "" || lang == "" {
		miss := make([]string, 0, 2)
		if userID == "" {
			miss = append(miss, "user_id")
		}
		if lang == "" {
			miss = append(miss, "lang_id")
		}
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeMissingFields, fmt.Sprintf("missing required query param(s): %s", strings.Join(miss, ", ")))
		return
	}
	// validate language ID
	var ok2 bool
	for _, l := range store.LanguagesList() {
		if l.ID == lang {
			ok2 = true
			break
		}
	}
	if !ok2 {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidLanguage, fmt.Sprintf("unsupported language id: %q", lang))
		return
	}
	packsCSV := strings.TrimSpace(q.Get("pack_ids"))
	var packs []string
	if packsCSV != "" {
		for _, p := range strings.Split(packsCSV, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				packs = append(packs, p)
			}
		}
		// validate packs exist by ID
		for _, p := range packs {
			if _, ok := store.GetPackByID(p); !ok {
				utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidPacks, fmt.Sprintf("unknown pack id: %s", p))
				return
			}
		}
	}
	// parse limit
	limit := 20
	if s := strings.TrimSpace(q.Get("limit")); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			if n > 100 {
				n = 100
			}
			limit = n
		}
	}

	vocabs := store.ListVocabs(userID, lang, packs)
	if len(vocabs) == 0 {
		utils.WriteOKData(w, []Flashcard{}, nil)
		return
	}
	// shuffle stable with seed
	rsrc := rand.New(rand.NewSource(time.Now().UnixNano()))
	rsrc.Shuffle(len(vocabs), func(i, j int) { vocabs[i], vocabs[j] = vocabs[j], vocabs[i] })
	if len(vocabs) > limit {
		vocabs = vocabs[:limit]
	}
	cards := make([]Flashcard, len(vocabs))
	for i, v := range vocabs {
		pack, _ := store.GetPackByID(v.PackID)
		cards[i] = Flashcard{ID: v.ID, Image: v.Image, Name: v.Name, PackName: pack.Name}
	}
	utils.WriteOKData(w, cards, map[string]any{"count": len(cards)})
}

// removed sanitizeFileBase: now lives in utils.UploadImage
