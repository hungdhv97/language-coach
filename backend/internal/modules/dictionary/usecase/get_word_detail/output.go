package get_word_detail

import (
	"github.com/english-coach/backend/internal/modules/dictionary/domain"
)

// Output represents detailed information about a word
type Output struct {
	Word           *domain.Word            `json:"word"`
	Senses         []SenseDetail           `json:"senses"`
	Pronunciations []*domain.Pronunciation `json:"pronunciations"`
	Relations      []*domain.WordRelation  `json:"relations,omitempty"`
}

// SenseDetail represents detailed information about a sense
type SenseDetail struct {
	ID                   int64             `json:"id"`
	SenseOrder           int16             `json:"sense_order"`
	PartOfSpeechID       int16             `json:"part_of_speech_id"`
	PartOfSpeechName     *string           `json:"part_of_speech_name,omitempty"`
	Definition           string            `json:"definition"`
	DefinitionLanguageID int16             `json:"definition_language_id"`
	LevelID              *int64            `json:"level_id,omitempty"`
	LevelName            *string           `json:"level_name,omitempty"`
	Note                 *string           `json:"note,omitempty"`
	Translations         []*domain.Word    `json:"translations"`
	Examples             []*domain.Example `json:"examples"`
}

