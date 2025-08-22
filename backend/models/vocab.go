package models

type Vocab struct {
	ID          string `json:"id"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	PackName    string `json:"pack_name"` // Foreign key to Pack.Name
	LangCode    string `json:"lang_code"` // redundant but useful for queries
}
