package model

// Language represents a supported language in the system
type Language struct {
	ID   int16  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
