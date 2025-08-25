package handlers

import (
	"net/http"

	"learnlang-backend/store"
	"learnlang-backend/utils"
)

// GetLanguagesHandler returns the list of supported languages.
func GetLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteOKData(w, store.LanguagesList(), nil)
}
