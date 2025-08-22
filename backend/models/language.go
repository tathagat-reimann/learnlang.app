package models

type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // e.g. "hi" for Hindi
}
