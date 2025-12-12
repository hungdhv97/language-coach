package model

// PartOfSpeech represents a grammatical category (noun, verb, adjective, etc.)
type PartOfSpeech struct {
	ID   int16  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
