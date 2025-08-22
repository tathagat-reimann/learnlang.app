package handlers

import (
	"encoding/json"
	"net/http"

	"learnlang-backend/models"
	"learnlang-backend/store"
	"learnlang-backend/utils"

	"github.com/google/uuid"
)

func CreatePackHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Pack
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.LangCode == "" || req.UserID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if language exists
	found := false
	for _, lang := range store.Languages {
		if lang.Code == req.LangCode {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Invalid language code", http.StatusBadRequest)
		return
	}

	// Check uniqueness
	key := utils.MakePackKey(req.UserID, req.LangCode, req.Name)
	if _, exists := store.PackIndex[key]; exists {
		http.Error(w, "Pack name already exists for this user and language", http.StatusConflict)
		return
	}

	// Create and store pack
	req.ID = uuid.New().String()
	store.Packs[req.ID] = req
	store.PackIndex[key] = req.ID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
