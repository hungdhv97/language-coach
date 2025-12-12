package model

// WordRelation represents a relationship between words
type WordRelation struct {
	RelationType string  `json:"relation_type"` // 'synonym', 'antonym', 'related'
	Note         *string `json:"note,omitempty"`
	TargetWord   *Word   `json:"target_word"`
}
