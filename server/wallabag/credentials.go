package wallabag

import (
	"time"
)

type WallabagOAuthToken struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    int     `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	Scope        *string `json:"scope"`
}

type Token struct {
	AccessToken  string
	ExpiresAt    time.Time
	RefreshToken string
}

func NewTokenFromPayload(payload *WallabagOAuthToken) *Token {
	expiresIn := time.Duration(payload.ExpiresIn) * time.Second
	deadline := time.Now().Add(expiresIn)
	return &Token{
		AccessToken:  payload.AccessToken,
		ExpiresAt:    deadline,
		RefreshToken: payload.RefreshToken,
	}
}

func (t *Token) HasExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// Username and password should not be stored in the database but the current implementation of oauth2 in wallabag
// has issues witht the refresh token. So we need to store the username and password to be able to refresh the token.
type Credentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	Token        *Token
}
