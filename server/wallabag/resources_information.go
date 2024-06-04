package wallabag

import (
	"encoding/json"
	"net/http"
)

type Information struct {
	*WallabagClient
}

type Info struct {
	AppName             string `json:"appname"`
	Version             string `json:"version"`
	AllowedRegistration bool   `json:"allowed_registration"`
}

func (c Information) Get() (*Info, error) {
	URL, _ := c.BuildURL(endpointInfo, nil)
	resp, err := c.CallAPI(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
