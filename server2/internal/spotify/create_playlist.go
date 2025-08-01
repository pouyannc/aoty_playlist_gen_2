package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreatePlaylist(trackURIs []string, token, uid, playlistName string) (string, error) {
	reqURL := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", uid)
	payload := struct {
		Name string `json:"name"`
	}{
		Name: playlistName,
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var newPlaylist struct {
		ID string `json:"id"`
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&newPlaylist)
	if err != nil {
		return "", err
	}

	fmt.Println("Created playlist ===================")

	err = PopulatePlaylist(trackURIs, newPlaylist.ID, token)
	if err != nil {
		return "", err
	}

	fmt.Println("Populated playlist ===================")

	return newPlaylist.ID, nil
}
