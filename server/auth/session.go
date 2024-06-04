package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecret() []byte {
	return []byte(os.Getenv("SESSION_SECRET"))
}

type SessionToken struct {
	Token    *jwt.Token
	IssuedAt time.Time
	UserID   string
}

func newSessionToken(userID string) (*SessionToken, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  userID,
	})

	return &SessionToken{
		Token:    token,
		IssuedAt: now,
		UserID:   userID,
	}, nil
}

func ParseSession(key string) (*SessionToken, error) {
	token, err := jwt.ParseWithClaims(key, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		sterr := fmt.Errorf("unexpected claims type: %T", token.Claims)
		return nil, &InvalidSessionError{key, sterr}
	}

	return &SessionToken{
		Token:    token,
		IssuedAt: claims.IssuedAt.Time,
		UserID:   claims.Subject,
	}, nil
}

func (t *SessionToken) Session() (string, error) {
	if t.Token.Raw != "" {
		return t.Token.Raw, nil
	}
	return t.Token.SignedString(getSecret())
}

func NewSessionKey(userID string) (string, error) {
	token, err := newSessionToken(userID)
	if err != nil {
		return "", err
	}
	return token.Session()
}
