package models

type Pack struct {
	ID     string `json:"id"`
	Name   string `json:"name"`    // Unique per user per language
	LangID string `json:"lang_id"` // foreign key to language.ID
	UserID string `json:"user_id"` // Optional for MVP
	Public bool   `json:"public"`  // Whether the pack is public or private
}
