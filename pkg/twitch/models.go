package twitch

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

type StreamInfo struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}
