package model

// Example represents an example sentence illustrating word usage
type Example struct {
	ID            int64                      `json:"id"`
	SourceSenseID int64                      `json:"source_sense_id"`
	LanguageID    int16                      `json:"language_id"`
	Content       string                     `json:"content"`
	AudioURL      *string                    `json:"audio_url,omitempty"`
	Source        *string                    `json:"source,omitempty"`
	Translations  []ExampleTranslationSimple `json:"translations,omitempty"`
}

// ExampleTranslationSimple represents a simple translation with language code
type ExampleTranslationSimple struct {
	Language string `json:"language"`
	Content  string `json:"content"`
}

// ExampleTranslation represents a translation of an example sentence
type ExampleTranslation struct {
	ID         int64  `json:"id"`
	ExampleID  int64  `json:"example_id"`
	LanguageID int16  `json:"language_id"`
	Content    string `json:"content"`
}
