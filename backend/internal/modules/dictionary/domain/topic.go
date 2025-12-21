package domain

// Topic represents a thematic category for organizing vocabulary
type Topic struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

