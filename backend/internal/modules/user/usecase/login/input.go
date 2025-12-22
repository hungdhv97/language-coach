package login

// Input represents the input for user login
type Input struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password string  `json:"password"`
}

