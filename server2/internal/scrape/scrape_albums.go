package scrape

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
)

type Album struct {
	Title  string
	Artist string
}

func ScrapeAlbums(page *rod.Page, scrapeURLs []string, filter string, nAlbums int) ([]*Album, error) {
	nAlbumsBuffer := 3
	nPages := len(scrapeURLs)
	totalAppended := 0
	albums := make([]*Album, (nAlbums+nAlbumsBuffer)*nPages)

	nErr := 0
	for i, u := range scrapeURLs {
		err := page.Navigate(u)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Current number of pages running in browser: %v\n", len(page.Browser().MustPages()))
		log.Printf("Page html: %v\n", page.MustHTML()[:400])
		err = page.Timeout(2000*time.Millisecond).WaitElementsMoreThan(".albumBlock", 0)
		if err != nil {
			nErr++
			fmt.Println("Elements failed to load on page")
			continue
		}

		albumElements, err := page.Timeout(2000 * time.Millisecond).Elements(".albumBlock")
		if err != nil {
			fmt.Println(err)
			continue
		}

		appendedFromPage := 0
		for j, e := range albumElements {
			albumSlicePosition := i + j*nPages
			if albumSlicePosition >= len(albums) {
				break
			}
			albums[albumSlicePosition] = &Album{
				Title:  e.MustElement(".albumTitle").MustText(),
				Artist: e.MustElement(".artistTitle").MustText(),
			}
			totalAppended++
			appendedFromPage++
		}

		fmt.Println("========= appended", appendedFromPage, "albums from page")
	}
	fmt.Printf("%v/%v pages loaded elements successfully\n", len(scrapeURLs)-nErr, len(scrapeURLs))
	if nErr >= len(scrapeURLs) {
		return []*Album{}, fmt.Errorf("failed to load album block elements from all pages")
	}

	fmt.Println("========= finished compiling albums slice")
	fmt.Println("scraped ", totalAppended, "albums")
	return albums, nil
}
