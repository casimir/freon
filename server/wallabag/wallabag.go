package wallabag

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	UserAgent = "freon"

	endpointToken = "/oauth/v2/token"
	endpointInfo  = "/api/info"
)

type Options interface {
	Validate() error
	ToMap() map[string]string
}

type WallabagClient struct {
	inner       http.Client
	credentials Credentials

	Entries     *Entries
	Information *Information
}

func NewWallabagClient(credentials Credentials) *WallabagClient {
	credentials.ServerURL = strings.TrimRight(credentials.ServerURL, "/")
	client := WallabagClient{
		credentials: credentials,
	}
	client.Entries = &Entries{&client}
	client.Information = &Information{&client}
	return &client
}

func (c WallabagClient) Token() *Token {
	return c.credentials.Token
}

func (c WallabagClient) BuildURL(path string, options Options) (string, error) {
	URL := c.credentials.ServerURL + path
	queryProvided := strings.Contains(path, "?")
	if !queryProvided && !(options == nil || reflect.ValueOf(options).IsNil()) {
		if err := options.Validate(); err != nil {
			return "", err
		}
		q := url.Values{}
		for k, v := range options.ToMap() {
			q.Set(k, v)
		}
		URL += "?" + q.Encode()
	}
	return URL, nil
}

func (c *WallabagClient) CallAPI(method string, URL string, payload any) (*http.Response, error) {
	needsAuth := !(strings.HasSuffix(URL, endpointToken) || strings.HasSuffix(URL, endpointInfo))

	if needsAuth {
		if c.credentials.Token == nil {
			return nil, &WallabagNotAuthenticatedError{}
		}
		if c.credentials.Token == nil || c.credentials.Token.HasExpired() {
			log.Printf("wallabag token has expired: %s", c.credentials.Token.ExpiresAt)
			// See documention of WallabagCredentials in server/auth/models.go
			// if err := c.RefreshToken(); err != nil {
			if err := c.FetchToken(c.credentials.Username, c.credentials.Password); err != nil {
				return nil, err
			}
		}
	}

	var req *http.Request
	var err error
	if payload != nil {
		body := new(bytes.Buffer)
		if err := json.NewEncoder(body).Encode(payload); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, URL, body)
	} else {
		req, err = http.NewRequest(method, URL, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/json")
	if needsAuth {
		req.Header.Add("Authorization", "Bearer "+c.credentials.Token.AccessToken)
	}

	start := time.Now()
	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("wallabag -> %s %s %d (%s)", method, URL, resp.StatusCode, time.Since(start))
	// wallabag returns 200 for any kind of successful operation
	if resp.StatusCode != http.StatusOK {
		return nil, &WallabagApiError{response: resp}
	}
	return resp, nil
}

func (c *WallabagClient) authenticate(grantData map[string]string) error {
	payload := map[string]string{
		"client_id":     c.credentials.ClientID,
		"client_secret": c.credentials.ClientSecret,
	}
	for k, v := range grantData {
		payload[k] = v
	}

	URL, _ := c.BuildURL(endpointToken, nil)
	resp, err := c.CallAPI(http.MethodPost, URL, payload)
	if err != nil {
		return err
	}

	var oauthToken WallabagOAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&oauthToken); err != nil {
		return err
	}
	c.credentials.Token = NewTokenFromPayload(&oauthToken)
	return nil
}

func (c *WallabagClient) FetchToken(username string, password string) error {
	return c.authenticate(map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	})
}

func (c *WallabagClient) RefreshToken() error {
	return c.authenticate(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": c.credentials.Token.RefreshToken,
	})
}
