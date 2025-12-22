package update_profile

// Input represents the input for updating user profile
type Input struct {
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BirthDay    *string `json:"birth_day,omitempty"` // Format: YYYY-MM-DD
	Bio         *string `json:"bio,omitempty"`
}

