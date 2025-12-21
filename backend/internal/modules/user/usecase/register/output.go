package register

// Output represents the output for user registration
type Output struct {
	UserID   int64   `json:"user_id"`
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
}

