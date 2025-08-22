package models

type Pack struct {
	ID       string `json:"id"`
	Name     string `json:"name"`      // Unique per user per language
	LangCode string `json:"lang_code"` // e.g. "hi"
	UserID   string `json:"user_id"`   // Optional for MVP
}
