package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// app matches the Mastadon app entitiy: https://docs.joinmastodon.org/api/entities/#app
type app struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Website      string `json:"website"`
	VapidKey     string `json:"vapid_key"`
}

// Authenticate runs through the full authentication flow
func Authenticate(instanceHost string) {
	newApp := createApp(instanceHost)
	fmt.Printf("%+v", newApp)
}

// createApp creates and returns an App struct: https://docs.joinmastodon.org/api/rest/apps/#post-api-v1-apps
func createApp(instanceHost string) app {
	requestURL, err := url.Parse(instanceHost + "/api/v1/apps")
	if err != nil {
		panic(err)
	}
	requestURL.Scheme = "https"

	appParams := map[string]string{"scopes": "write follow read",
		"redirect_uris": "urn:ietf:wg:oauth:2.0:oob", "client_name": "fedapp"}

	requestQuery, err := json.Marshal(appParams)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(requestURL.String(), "application/json", bytes.NewBuffer(requestQuery))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	var newApp app
	err = decoder.Decode(&newApp)
	if err != nil {
		panic(err)
	}

	return newApp
}