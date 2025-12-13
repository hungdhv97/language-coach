package model

import "time"

// Word represents a vocabulary word in a specific language
type Word struct {
	ID              int64     `json:"id"`
	LanguageID      int16     `json:"language_id"`
	Lemma           string    `json:"lemma"`
	LemmaNormalized *string   `json:"lemma_normalized,omitempty"`
	SearchKey       *string   `json:"search_key,omitempty"`
	Romanization    *string   `json:"romanization,omitempty"`
	ScriptCode      *string   `json:"script_code,omitempty"`
	FrequencyRank   *int      `json:"frequency_rank,omitempty"`
	Note            *string   `json:"note,omitempty"`
	Topics          []*Topic  `json:"topics,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
