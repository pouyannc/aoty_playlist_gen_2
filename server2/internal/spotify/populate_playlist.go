package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PopulatePlaylist(trackURIs []string, playlistID, token string) error {
	fmt.Println("playlistID: ", playlistID)

	reqURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
	payload := struct {
		URIs []string `json:"uris"`
	}{
		URIs: trackURIs,
	}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("spotify returned status %d, and error reading body: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("spotify returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
