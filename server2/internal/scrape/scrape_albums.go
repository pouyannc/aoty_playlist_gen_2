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
	startTime := time.Now()
	nPages := len(albumsPages)
	albumElements := make([]rod.Elements, nPages)
	totalElements := 0

	nErr := 0
	for i, page := range albumsPages {
		err := page.Timeout(800*time.Millisecond).WaitElementsMoreThan(".albumBlock", 0)
		if err != nil {
			nErr++
			fmt.Println("Elements failed to load on page")
			continue
		}

		fmt.Println("========= loaded enough elements in : ", time.Since(startTime))
		startTime = time.Now()

		albumElements[i] = page.MustElements(".albumBlock")
		totalElements += len(albumElements[i])

		fmt.Println("========= appended", len(albumElements[i]), "elements in : ", time.Since(startTime))
		startTime = time.Now()
	}
	if nErr >= len(albumsPages) {
		return []Album{}, fmt.Errorf("failed to load album block elements from all pages")
	}

	fmt.Println("========= finished collecting all elements in : ", time.Since(startTime))
	startTime = time.Now()

	albums := []Album{}
	var wg sync.WaitGroup

	fmt.Println("Got ", totalElements, "album elements")

	for i := 0; i < nrAlbums+3 && i < totalElements; i++ {
		nPage := i % len(albumElements)
		nAlbum := i / len(albumElements)

		// need this check since n album elements on each page are not equal lengths
		if nAlbum >= len(albumElements[nPage]) {
			nrAlbums++
			totalElements++
			continue
		}

		wg.Add(1)
		go func(e *rod.Element) {
			defer wg.Done()
			albums = append(albums, Album{
				Title:  e.MustElement(".albumTitle").MustText(),
				Artist: e.MustElement(".artistTitle").MustText(),
			})
		}(albumElements[nPage][nAlbum])
	}

	wg.Wait()

	fmt.Println("========= finished compiling albums slice in : ", time.Since(startTime))
	fmt.Println("scraped ", len(albums), "albums: ", albums)
	return albums, nil
}
