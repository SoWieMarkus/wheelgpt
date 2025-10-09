package twitch

// User represents a Twitch user with relevant attributes.
type User struct {
	// ID of the user
	ID string
	// Login name of the user
	Name string
	// Display name of the user
	DisplayName string
	// Chat color of the user
	Color string
	// Roles of the user in the channel
	IsMod         bool
	IsBroadcaster bool
	IsSubscriber  bool
	IsVIP         bool
	// Badges of the user
	Badges map[string]int
}
