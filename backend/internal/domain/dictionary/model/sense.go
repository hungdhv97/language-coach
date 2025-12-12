package model

// Sense represents a specific meaning or definition of a word
type Sense struct {
	ID                   int64   `json:"id"`
	WordID               int64   `json:"word_id"`
	SenseOrder           int16   `json:"sense_order"`
	PartOfSpeechID       int16   `json:"part_of_speech_id"`
	Definition           string  `json:"definition"`
	DefinitionLanguageID int16   `json:"definition_language_id"`
	UsageLabel           *string `json:"usage_label,omitempty"`
	LevelID              *int64  `json:"level_id,omitempty"`
	Note                 *string `json:"note,omitempty"`
}
