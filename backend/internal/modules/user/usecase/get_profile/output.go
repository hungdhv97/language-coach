package get_profile

// Output represents the output for getting user profile
type Output struct {
	UserID      int64   `json:"user_id"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BirthDay    *string `json:"birth_day,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

