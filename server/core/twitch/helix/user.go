package helix

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	CreatedAt       string `json:"created_at"`
}

// See: https://dev.twitch.tv/docs/api/reference#get-users
func (c *Client) GetUsers() ([]User, error) {
	var users []User
	_, err := c.get("/users", &users, nil)
	if err != nil {
		return nil, err
	}
	return users, nil
}
