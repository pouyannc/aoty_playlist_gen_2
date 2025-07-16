package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetUID(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("spotify API error: %s", string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse user info: %v", err)
	}

	return result.ID, nil
}
