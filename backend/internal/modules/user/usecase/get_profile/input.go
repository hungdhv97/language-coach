package get_profile

// Input represents the input for getting user profile
type Input struct {
	UserID int64 `json:"user_id" validate:"required,gt=0"`
}

