package twitch

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// Scopes requested in the OAuth flows.
var (
	// StreamerScopes is all a streamer grants: the bot may read and write
	// chat in their channel. Identity needs no scope.
	StreamerScopes = []string{"channel:bot"}
	// BotScopes is what the bot account grants once.
	BotScopes = []string{"user:read:chat", "user:write:chat", "user:bot"}
)

// AuthorizeURL builds the Twitch authorization URL.
func (c *Client) AuthorizeURL(redirectURI string, scopes []string, state string, forceVerify bool) string {
	q := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
		"redirect_uri":  {redirectURI},
		"scope":         {strings.Join(scopes, " ")},
		"state":         {state},
	}
	if forceVerify {
		q.Set("force_verify", "true")
	}
	return c.AuthURL + "/authorize?" + q.Encode()
}

// UserToken is the result of a code exchange.
type UserToken struct {
	AccessToken string   `json:"access_token"`
	Scopes      []string `json:"scope"`
}

// ExchangeCode trades an authorization code for a user token. The token is
// only used once to identify the user and then dropped.
func (c *Client) ExchangeCode(ctx context.Context, code, redirectURI string) (UserToken, error) {
	form := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
	}
	var t UserToken
	err := c.postForm(ctx, "user_token", c.AuthURL+"/token", form, &t)
	return t, err
}

// RevokeToken revokes a user token (best effort).
func (c *Client) RevokeToken(ctx context.Context, token string) {
	form := url.Values{"client_id": {c.ClientID}, "token": {token}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AuthURL+"/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_ = c.do(req, "revoke", nil)
}
