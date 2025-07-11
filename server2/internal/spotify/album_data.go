package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pouyannc/aoty_list_gen/internal/scrape"
)

type SpotifySearchResp struct {
	Albums struct {
		Items []struct {
			ID     string `json:"id"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"albums"`
}

type SpotifyAlbum struct {
	AlbumID  string
	CoverURL string
	Artist   string
}

func AlbumData(albums []scrape.Album, token string) ([]SpotifyAlbum, error) {

	baseURL := "https://api.spotify.com/v1/search"
	params := url.Values{}
	params.Set("type", "album")
	params.Set("limit", "2")

	accessToken := token

	spotifyAlbumData := []SpotifyAlbum{}

	for _, album := range albums {
		params.Set("q", fmt.Sprintf("%s %s", album.Title, album.Artist))
		searchURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		req, err := http.NewRequest("GET", searchURL, nil)
		if err != nil {
			return []SpotifyAlbum{}, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return []SpotifyAlbum{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return []SpotifyAlbum{}, errors.New(string(body))
		}

		result := SpotifySearchResp{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return []SpotifyAlbum{}, err
		}

		for _, searchItem := range result.Albums.Items {
			if searchItem.Name == album.Title && searchItem.Artists[0].Name == album.Artist {
				spotifyAlbumData = append(spotifyAlbumData, SpotifyAlbum{
					AlbumID:  searchItem.ID,
					CoverURL: searchItem.Images[1].URL,
					Artist:   album.Artist,
				})
				break
			}
		}
	}

	return spotifyAlbumData, nil
}
