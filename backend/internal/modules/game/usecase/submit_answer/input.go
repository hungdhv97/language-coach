package submit_answer

// Input represents the input to submit an answer
type Input struct {
	QuestionID       int64 `json:"question_id" validate:"required,gt=0"`
	SelectedOptionID int64 `json:"selected_option_id" validate:"required,gt=0"`
	ResponseTimeMs   *int  `json:"response_time_ms,omitempty" validate:"omitempty,gt=0"`
}

