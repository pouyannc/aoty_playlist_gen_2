package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

type AlbumScrapeData interface {
	GetTitleArtist() (string, string)
}

func AlbumData[T AlbumScrapeData](albums []*T, token string, nMaxAlbums int) ([]SpotifyAlbum, error) {
	nAlbumSearchBuffer := 3

	baseURL := "https://api.spotify.com/v1/search"

	spotifyAlbumData := []SpotifyAlbum{}

	nAlbums := 0

	//var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for _, album := range albums {
		if nAlbums >= nMaxAlbums+nAlbumSearchBuffer {
			break
		}
		if album == nil {
			continue
		}
		nAlbums++
		wg.Add(1)
		sem <- struct{}{}

		go func(alb T) {
			defer wg.Done()
			defer func() { <-sem }()

			albumTitle, albumArtist := alb.GetTitleArtist()

			q := url.QueryEscape(albumTitle + " " + albumArtist)
			searchURL := fmt.Sprintf(
				"%s?type=album&limit=2&q=%s",
				baseURL,
				q,
			)

			req, err := http.NewRequest("GET", searchURL, nil)
			if err != nil {
				log.Printf("error creating request: %v\n", err)
				return
			}
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("error getting response: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				log.Printf("response status with bad status code: %v\n", string(body))
				return
			}

			var result SpotifySearchResp
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				log.Printf("error deconding json: %v\n", err)
				return
			}

			for _, searchItem := range result.Albums.Items {
				titleCompare := strings.EqualFold(searchItem.Name[0:1], albumTitle[0:1]) &&
					strings.EqualFold(searchItem.Name[len(searchItem.Name)-1:], albumTitle[len(albumTitle)-1:])
				artistCompare := strings.EqualFold(searchItem.Artists[0].Name[0:1], albumArtist[0:1]) &&
					strings.EqualFold(searchItem.Artists[0].Name[len(searchItem.Artists[0].Name)-1:], albumArtist[len(albumArtist)-1:])

				if titleCompare && artistCompare {
					//mu.Lock()
					spotifyAlbumData = append(spotifyAlbumData, SpotifyAlbum{
						AlbumID:  searchItem.ID,
						CoverURL: searchItem.Images[1].URL,
						Artist:   albumArtist,
					})
					//mu.Unlock()
					break
				}
			}
		}(*album)
	}

	wg.Wait()
	return spotifyAlbumData, nil
}
