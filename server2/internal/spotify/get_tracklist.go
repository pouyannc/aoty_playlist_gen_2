package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AlbumsResp struct {
	Albums []struct {
		Tracks struct {
			Items []struct {
				URI string `json:"uri"`
			} `json:"items"`
		} `json:"tracks"`
	} `json:"albums"`
}

func GetTracklist(albumData []SpotifyAlbum, tracksPerAlbum, nTracks int, token string) ([]string, error) {
	albumIDs := []string{}
	for _, album := range albumData {
		albumIDs = append(albumIDs, album.AlbumID)
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/albums?ids=0", nil)
	if err != nil {
		return []string{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println(albumIDs, len(albumIDs))

	tracklist := []string{}
	maxGetAlbumAmount := 20
	for i := 0; i < len(albumIDs); i += maxGetAlbumAmount {
		query := req.URL.Query()
		var q string
		if len(albumIDs)-i < maxGetAlbumAmount {
			q = strings.Join(albumIDs[i:], ",")
		} else {
			q = strings.Join(albumIDs[i:i+maxGetAlbumAmount], ",")
		}
		query.Set("ids", q)
		req.URL.RawQuery = query.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()

		var albumTracks AlbumsResp
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&albumTracks)
		if err != nil {
			return []string{}, err
		}

		for _, album := range albumTracks.Albums {
			randomTracks, err := pickNRandomTracks([]struct{ URI string }(album.Tracks.Items), tracksPerAlbum)
			if err != nil {
				return []string{}, err
			}
			tracklist = append(tracklist, randomTracks...)
		}

		if len(tracklist) >= nTracks {
			break
		}
	}

	return tracklist[:nTracks], nil
}
