package events

// See: https://dev.twitch.tv/docs/eventsub/eventsub-reference/#channel-chat-message-event
type ChatMessageEvent struct {
	BroadcasterUserID           string      `json:"broadcaster_user_id"`
	BroadcasterUserName         string      `json:"broadcaster_user_name"`
	BroadcasterUserLogin        string      `json:"broadcaster_user_login"`
	ChatterUserID               string      `json:"chatter_user_id"`
	ChatterUserName             string      `json:"chatter_user_name"`
	ChatterUserLogin            string      `json:"chatter_user_login"`
	MessageID                   string      `json:"message_id"`
	Message                     ChatMessage `json:"message"`
	MessageType                 string      `json:"message_type"`
	Badges                      []Badge     `json:"badges"`
	Cheer                       *CheerInfo  `json:"cheer,omitempty"`
	Color                       string      `json:"color"`
	Reply                       *ReplyInfo  `json:"reply,omitempty"`
	ChannelPointsCustomRewardID string      `json:"channel_points_custom_reward_id,omitempty"`
	SourceBroadcasterUserID     string      `json:"source_broadcaster_user_id,omitempty"`
	SourceBroadcasterUserName   string      `json:"source_broadcaster_user_name,omitempty"`
	SourceBroadcasterUserLogin  string      `json:"source_broadcaster_user_login,omitempty"`
	SourceMessageID             string      `json:"source_message_id,omitempty"`
	SourceBadges                []Badge     `json:"source_badges,omitempty"`
	IsSourceOnly                bool        `json:"is_source_only,omitempty"`
}

type ChatMessage struct {
	Text      string                `json:"text"`
	Fragments []ChatMessageFragment `json:"fragments"`
}

type ChatMessageFragment struct {
	Type      string         `json:"type"`
	Text      string         `json:"text"`
	Cheermote *CheermoteInfo `json:"cheermote,omitempty"`
	Emote     *EmoteInfo     `json:"emote,omitempty"`
	Mention   *MentionInfo   `json:"mention,omitempty"`
}

type CheermoteInfo struct {
	Prefix string `json:"prefix"`
	Bits   int    `json:"bits"`
	Tier   int    `json:"tier"`
}

type EmoteInfo struct {
	ID         string   `json:"id"`
	EmoteSetID string   `json:"emote_set_id"`
	OwnerID    string   `json:"owner_id"`
	Format     []string `json:"format"`
}

type MentionInfo struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type Badge struct {
	SetID string `json:"set_id"`
	ID    string `json:"id"`
	Info  string `json:"info,omitempty"`
}

type CheerInfo struct {
	Bits int `json:"bits"`
}

type ReplyInfo struct {
	ParentMessageID   string `json:"parent_message_id"`
	ParentMessageBody string `json:"parent_message_body"`
	ParentUserID      string `json:"parent_user_id"`
	ParentUserName    string `json:"parent_user_name"`
	ParentUserLogin   string `json:"parent_user_login"`
	ThreadMessageID   string `json:"thread_message_id"`
}
