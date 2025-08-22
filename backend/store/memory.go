package store

import (
	"learnlang-backend/models"
)

var Languages = []models.Language{
	{ID: "1", Name: "Hindi", Code: "hi"},
	{ID: "2", Name: "German", Code: "de"},
	// Add more as needed
}

// Pack store: key is Pack ID
var Packs = map[string]models.Pack{}

// Vocab store: key is Vocab ID
var VocabEntries = map[string]models.Vocab{}

// Composite index to enforce uniqueness of Pack name per user per language
// Key format: userID:langCode:packName → packID
var PackIndex = map[string]string{}
