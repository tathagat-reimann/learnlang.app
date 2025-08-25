package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"learnlang-backend/models"
	"learnlang-backend/store"
	"learnlang-backend/utils"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
)

// CreatePackRequestDTO is the request DTO for creating a new pack.
// It intentionally omits ID to prevent clients from setting it.
type CreatePackRequestDTO struct {
	Name   string `json:"name"`
	LangID string `json:"lang_id"`
	UserID string `json:"user_id"`
}

// GetPacksHandler returns all packs in the in-memory store.
func GetPacksHandler(w http.ResponseWriter, r *http.Request) {
	packs := store.GetAllPacks()
	utils.WriteOKData(w, packs, nil)
}

func CreatePackHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePackRequestDTO
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF):
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeEmptyBody, "request body must not be empty")
		case errors.As(err, &syntaxErr):
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeJSONSyntax, fmt.Sprintf("badly-formed JSON at position %d", syntaxErr.Offset))
		case errors.As(err, &typeErr):
			if typeErr.Field != "" {
				utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeJSONType, fmt.Sprintf("invalid type for field %q: expected %s", typeErr.Field, typeErr.Type.String()))
			} else {
				utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeJSONType, fmt.Sprintf("invalid JSON value at position %d: expected %s", typeErr.Offset, typeErr.Type.String()))
			}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeUnknownField, fmt.Sprintf("unknown field %s", field))
		case errors.Is(err, io.ErrUnexpectedEOF):
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeJSONSyntax, "badly-formed JSON")
		default:
			utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidJSON, "invalid JSON payload")
		}
		return
	}

	// Body must contain a single JSON object
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeMultipleObjects, "request body must contain a single JSON object")
		return
	}

	// Basic validation
	missing := make([]string, 0, 3)
	if req.Name == "" {
		missing = append(missing, "name")
	}
	if req.LangID == "" {
		missing = append(missing, "lang_id")
	}
	if req.UserID == "" {
		missing = append(missing, "user_id")
	}
	if len(missing) > 0 {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeMissingFields, fmt.Sprintf("missing required field(s): %s", strings.Join(missing, ", ")))
		return
	}

	// Check if language exists and get its ID
	// Validate language by ID
	langs := store.LanguagesList()
	var found bool
	for _, l := range langs {
		if l.ID == req.LangID {
			found = true
			break
		}
	}
	if !found {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidLanguage, fmt.Sprintf("unsupported language id: %q", req.LangID))
		return
	}

	// Check uniqueness
	key := utils.MakePackKey(req.UserID, req.LangID, req.Name)
	if store.PackExistsByKey(key) {
		utils.WriteErrorWithRequest(w, r, http.StatusConflict, utils.CodeDuplicatePack, fmt.Sprintf("pack %q already exists for user %q and language %q", req.Name, req.UserID, req.LangID))
		return
	}

	// Create and store pack
	pack := models.Pack{
		ID:     uuid.New().String(),
		Name:   req.Name,
		LangID: req.LangID,
		UserID: req.UserID,
	}
	store.CreatePack(pack, key)

	utils.WriteCreatedData(w, pack, nil)
}

// GetPackByIDHandler returns a single pack by its ID.
func GetPackByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		utils.WriteErrorWithRequest(w, r, http.StatusBadRequest, utils.CodeInvalidPack, "missing pack id")
		return
	}
	p, ok := store.GetPackByID(id)
	if !ok {
		utils.WriteErrorWithRequest(w, r, http.StatusNotFound, utils.CodeInvalidPack, fmt.Sprintf("unknown pack id: %q", id))
		return
	}
	// Fetch related vocabs for this pack (user/lang implied by pack)
	vocabs := store.ListVocabsByPackID(p.ID)
	type response struct {
		Pack   models.Pack    `json:"pack"`
		Vocabs []models.Vocab `json:"vocabs"`
	}
	utils.WriteOKData(w, response{Pack: p, Vocabs: vocabs}, nil)
}
