package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type webClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

type pluginClaims struct {
	jwt.RegisteredClaims
	ID    string `json:"id"`
	Token string `json:"token"`
}

func IssueWebToken(channelID, secret string) (string, error) {
	claims := webClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		ID: channelID,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func VerifyWebToken(tokenStr, secret string) (channelID string, err error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &webClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	c, ok := tok.Claims.(*webClaims)
	if !ok || !tok.Valid {
		return "", errors.New("invalid token")
	}
	return c.ID, nil
}

func IssuePluginToken(channelID, channelToken, secret string) (string, error) {
	claims := pluginClaims{
		ID:    channelID,
		Token: channelToken,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func VerifyPluginToken(tokenStr, secret string) (channelID, channelToken string, err error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &pluginClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", "", err
	}
	c, ok := tok.Claims.(*pluginClaims)
	if !ok || !tok.Valid {
		return "", "", errors.New("invalid token")
	}
	return c.ID, c.Token, nil
}
