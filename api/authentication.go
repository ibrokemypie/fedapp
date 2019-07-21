package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// App matches the Mastadon app entitiy: https://docs.joinmastodon.org/api/entities/#app
type App struct {
	clientID     string
	clientSecret string
	clientName   string
	redirectUris string
	scopes       string
}

// CreateApp creates and returns an App struct
func CreateApp(instanceHost string) App {
	newApp := App{scopes: "write follow read", redirectUris: "urn:ietf:wg:oauth:2.0:oob",
		clientName: "fedapp"}

	requestURL, err := url.Parse(instanceHost + "/api/v1/apps")
	if err != nil {
		panic(err)
	}
	requestURL.Scheme = "https"

	requestQuery, err := json.Marshal(map[string]string{"scopes": newApp.scopes,
		"redirect_uris": newApp.redirectUris, "client_name": newApp.clientName})
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(requestURL.String(), "application/json", bytes.NewBuffer(requestQuery))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return newApp
}
