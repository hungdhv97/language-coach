package get_word_detail

// Input represents the input for getting word detail
type Input struct {
	WordID int64 `json:"word_id" validate:"required,gt=0"`
}

