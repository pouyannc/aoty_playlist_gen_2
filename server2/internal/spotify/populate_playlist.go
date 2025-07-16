package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func PopulatePlaylist(trackURIs []string, playlistID, token string) (string, error) {
	reqURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
	payload := struct {
		URIs []string `json:"uris"`
	}{
		URIs: trackURIs,
	}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var playlistResp struct {
		SnapshotID string `json:"snapshot_id"`
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&playlistResp)
	if err != nil {
		return "", err
	}

	return playlistResp.SnapshotID, nil
}
