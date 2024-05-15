package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

var t *Token
var tokenNextExpiry = time.Now()

func refreshToken() error {
	req, err := http.NewRequest(http.MethodPost, settings.OAuth2.TokenEndpoint, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		log.Fatal(err)
	}

	authorization := base64.StdEncoding.EncodeToString([]byte(settings.OAuth2.ClientId + ":" + settings.OAuth2.ClientSecret))
	req.Header.Add("Authorization", "Basic "+authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	requestTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respString, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(respString, &t)
	if err != nil {
		return err
	}

	tokenNextExpiry = requestTime.Add(time.Second * time.Duration(t.ExpiresIn))
	return nil
}

func invalidateToken() {
	log.Print("invalidating token")
	t = nil
	tokenNextExpiry = time.Now()
}

func getToken() (*Token, error) {
	if t != nil && time.Now().Before(tokenNextExpiry) {
		log.Print("token still valid")
		return t, nil
	}

	log.Print("retrieving new token")
	err := refreshToken()
	if err != nil {
		return nil, err
	}

	return t, nil
}
