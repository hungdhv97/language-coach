package domain

// Pronunciation represents pronunciation information for a word
type Pronunciation struct {
	ID       int64   `json:"id"`
	WordID   int64   `json:"word_id"`
	Dialect  *string `json:"dialect,omitempty"`
	IPA      *string `json:"ipa,omitempty"`
	Phonetic *string `json:"phonetic,omitempty"`
	AudioURL *string `json:"audio_url,omitempty"`
}
