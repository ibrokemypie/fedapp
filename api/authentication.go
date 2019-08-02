package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/therecipe/qt/core"
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
func Authenticate(instanceHost string, authChan chan string, settings *core.QSettings) {
	newApp := createApp(instanceHost)
	// fmt.Printf("%+v", newApp)

	authURL := authorizeApp(instanceHost, newApp)
	authChan <- authURL

	authCode := <-authChan
	accesstoken := getAccessToken(instanceHost, newApp, authCode)
	fmt.Println(accesstoken)
	settings.SetValue("access_token", core.NewQVariant1(accesstoken))
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

// authorizeApp creates the url for the user to authorize the app
func authorizeApp(instanceHost string, newApp app) string {
	requestURL, err := url.Parse(instanceHost + "/oauth/authorize")
	if err != nil {
		panic(err)
	}
	requestURL.Scheme = "https"

	v := url.Values{}
	v.Set("scope", "write follow read")
	v.Set("response_type", "code")
	v.Set("redirect_uri", "urn:ietf:wg:oauth:2.0:oob")
	v.Set("client_id", newApp.ClientID)

	requestURL.RawQuery = v.Encode()

	return requestURL.String()
}

func getAccessToken(instanceHost string, newApp app, authCode string) string {
	requestURL, err := url.Parse(instanceHost + "/oauth/token")
	if err != nil {
		panic(err)
	}
	requestURL.Scheme = "https"

	v := url.Values{}
	v.Set("client_id", newApp.ClientID)
	v.Set("client_secret", newApp.ClientSecret)
	v.Set("grant_type", "authorization_code")
	v.Set("code", authCode)
	v.Set("redirect_uri", newApp.RedirectURI)

	requestURL.RawQuery = v.Encode()

	resp, err := http.Post(requestURL.String(), "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	type accessTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var newResponse accessTokenResponse

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&newResponse)
	if err != nil {
		panic(err)
	}

	return newResponse.AccessToken
}
