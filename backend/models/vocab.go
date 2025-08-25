package models

type Vocab struct {
	ID          string `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Translation string `json:"translation"`
	PackID      string `json:"pack_id"` // foreign key to Pack.ID
}
