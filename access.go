package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type tokenResponseJSON struct {
	AccessToken string `json:"access_token"`
}

func getAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, blizzardTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp tokenResponseJSON
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}
