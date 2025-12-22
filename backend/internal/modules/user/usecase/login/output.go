package login

// Output represents the output for user login
type Output struct {
	Token    string  `json:"token"`
	UserID   int64   `json:"user_id"`
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
}

