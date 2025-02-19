package veeam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Username string
	Password string
	Endpoint string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	// ...other fields if needed...
}

func (c *Config) Client() (interface{}, error) {
	url := c.Endpoint + "/api/oauth2/token"
	data := []byte(`grant_type=password&username=` + c.Username + `&password=` + c.Password)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-api-version", "1.2-rev0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication issue: received status code %d", resp.StatusCode)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return tokenResponse.AccessToken, nil
}
