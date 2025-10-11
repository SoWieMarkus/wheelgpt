package events

// See: https://dev.twitch.tv/docs/eventsub/eventsub-reference/#channel-cheer-event
type ChannelCheerEvent struct {
	IsAnonymous          bool   `json:"is_anonymous"`
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Message              string `json:"message"`
	Bits                 int    `json:"bits"`
}
