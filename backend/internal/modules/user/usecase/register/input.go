package register

// Input represents the input for user registration
type Input struct {
	DisplayName *string `json:"display_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Username    *string `json:"username,omitempty"`
	Password    string  `json:"password"`
}

