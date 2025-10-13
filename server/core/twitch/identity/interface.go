package identity

type TwitchIdentityAPI interface {
	RequestUserAccessToken(code, redirectURI string) (*UserAccessToken, error)
	RequestAppAccessToken() (*AppAccessToken, error)
}
