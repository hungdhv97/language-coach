package http

import (
	"time"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
)

// WordResponse represents a word for HTTP response
type WordResponse struct {
	ID              int64           `json:"id"`
	LanguageID      int16           `json:"language_id"`
	Lemma           string          `json:"lemma"`
	LemmaNormalized *string         `json:"lemma_normalized,omitempty"`
	SearchKey       *string         `json:"search_key,omitempty"`
	Romanization    *string         `json:"romanization,omitempty"`
	ScriptCode      *string         `json:"script_code,omitempty"`
	FrequencyRank   *int            `json:"frequency_rank,omitempty"`
	Note            *string         `json:"note,omitempty"`
	Topics          []*domain.Topic `json:"topics,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// mapWordToResponse maps domain.Word to WordResponse
func mapWordToResponse(word *domain.Word) *WordResponse {
	if word == nil {
		return nil
	}
	return &WordResponse{
		ID:              word.ID,
		LanguageID:      word.LanguageID,
		Lemma:           word.Lemma,
		LemmaNormalized: word.LemmaNormalized,
		SearchKey:       word.SearchKey,
		Romanization:    word.Romanization,
		ScriptCode:      word.ScriptCode,
		FrequencyRank:   word.FrequencyRank,
		Note:            word.Note,
		Topics:          word.Topics,
		CreatedAt:       word.CreatedAt,
		UpdatedAt:       word.UpdatedAt,
	}
}

// mapWordsToResponse maps slice of domain.Word to slice of WordResponse
func mapWordsToResponse(words []*domain.Word) []*WordResponse {
	if words == nil {
		return nil
	}
	result := make([]*WordResponse, len(words))
	for i, word := range words {
		result[i] = mapWordToResponse(word)
	}
	return result
}

// SearchWordsRequest represents the query parameters for word search
type SearchWordsRequest struct {
	Query      string `form:"q" binding:"required"`
	LanguageID int16  `form:"languageId" binding:"required"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
}

// GetLevelsRequest represents the query parameters for getting levels
type GetLevelsRequest struct {
	LanguageID *int16 `form:"languageId"`
}

// GetWordDetailRequest represents the path parameter for getting word detail
type GetWordDetailRequest struct {
	WordID int64 `uri:"wordId" binding:"required"`
}

// WordRelationResponse represents a word relation for HTTP response
type WordRelationResponse struct {
	RelationType string        `json:"relation_type"`
	Note         *string       `json:"note,omitempty"`
	TargetWord   *WordResponse `json:"target_word"`
}

// SenseDetailResponse represents detailed information about a sense for HTTP response.
type SenseDetailResponse struct {
	ID                   int64           `json:"id"`
	SenseOrder           int16           `json:"sense_order"`
	PartOfSpeechID       int16           `json:"part_of_speech_id"`
	PartOfSpeechName     *string         `json:"part_of_speech_name,omitempty"`
	Definition           string          `json:"definition"`
	DefinitionLanguageID int16           `json:"definition_language_id"`
	LevelID              *int64          `json:"level_id,omitempty"`
	LevelName            *string         `json:"level_name,omitempty"`
	Note                 *string         `json:"note,omitempty"`
	Translations         []*WordResponse `json:"translations,omitempty"`
	Examples             interface{}     `json:"examples,omitempty"`
}

// GetWordDetailResponse represents the HTTP response for getting word detail.
type GetWordDetailResponse struct {
	Word           *WordResponse           `json:"word"`
	Senses         []SenseDetailResponse   `json:"senses"`
	Pronunciations interface{}             `json:"pronunciations"`
	Relations      []*WordRelationResponse `json:"relations,omitempty"`
}
