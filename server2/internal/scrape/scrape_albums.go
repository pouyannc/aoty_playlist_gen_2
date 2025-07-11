package scrape

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

type Album struct {
	Title  string
	Artist string
}

func ScrapeAlbums(albumsPage *rod.Page, nrAlbums int) ([]Album, error) {
	fmt.Println("in scrape function")
	err := albumsPage.Timeout(3*time.Second).WaitElementsMoreThan(".albumBlock", 0)
	if err != nil {
		return []Album{}, err
	}

	fmt.Println("loaded elements.")

	albums := []Album{}
	albumElements := albumsPage.MustElements(".albumBlock")

	for _, ele := range albumElements {
		albums = append(albums, Album{
			Title:  ele.MustElement(".albumTitle").MustText(),
			Artist: ele.MustElement(".artistTitle").MustText(),
		})

		if len(albums) == nrAlbums {
			break
		}
	}

	return albums, nil
}
