package twitch

type Subscription struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition map[string]string `json:"condition"`
	Transport map[string]string `json:"transport"`
	CreatedAt string            `json:"created_at"`
	Cost      int               `json:"cost"`
}

type MessageFragment struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Text      string            `json:"text"`
	Fragments []MessageFragment `json:"fragments"`
}

type Event struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ChatterUserID        string `json:"chatter_user_id"`
	ChatterUserLogin     string `json:"chatter_user_login"`
	ChatterUserName      string `json:"chatter_user_name"`
	MessageId            string `json:"message_id"`
}

type ChatMessage struct {
	Subscription Subscription `json:"subscription"`
	Event        Event        `json:"event"`
}

func (message ChatMessage) GetUser() User {
	return User{}
}
