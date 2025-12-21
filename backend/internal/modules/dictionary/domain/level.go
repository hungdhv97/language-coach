package domain

// Level represents a difficulty or proficiency level (e.g., HSK1, A1, N3)
type Level struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	LanguageID      *int16  `json:"language_id,omitempty"`
	DifficultyOrder *int16  `json:"difficulty_order,omitempty"`
}

