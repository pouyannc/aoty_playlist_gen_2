package scrape

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type Album struct {
	Title  string
	Artist string
}

func ScrapeAlbums(albumsPages []*rod.Page, filter string, nrAlbums int) ([]Album, error) {
	fmt.Println("in scrape function")
	var albumElements []rod.Elements
	for _, page := range albumsPages {
		err := page.Timeout(3*time.Second).WaitElementsMoreThan(".albumBlock", 0)
		if err != nil {
			return []Album{}, err
		}

		fmt.Println("loaded elements.")

		albumElements = append(albumElements, page.MustElements(".albumBlock"))
	}
	albums := []Album{}

	fmt.Println("got elements")
	//var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < nrAlbums+3 && i < len(albumElements)*len(albumElements[0]); i++ {
		wg.Add(1)

		nPage := i % len(albumElements)
		nAlbum := i / len(albumElements)

		go func(e *rod.Element) {
			defer wg.Done()
			albums = append(albums, Album{
				Title:  e.MustElement(".albumTitle").MustText(),
				Artist: e.MustElement(".artistTitle").MustText(),
			})
		}(albumElements[nPage][nAlbum])
	}

	wg.Wait()
	fmt.Println("scraped albums")
	return albums, nil
}
