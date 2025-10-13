package identity

type IdentityAPI interface {
	RequestUserAccessToken(code, redirectURI string) (*UserAccessToken, error)
	RequestAppAccessToken() (*AppAccessToken, error)
}
