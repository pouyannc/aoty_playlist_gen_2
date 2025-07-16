package spotify

import (
	"fmt"
	"math/rand"
	"time"
)

func pickNRandomTracks(tracks []struct{ URI string }, n int) ([]string, error) {
	if n < 0 {
		return []string{}, fmt.Errorf("n cannot be negative")
	}
	if n > len(tracks) {
		return []string{}, fmt.Errorf("n cannot be greater than the length of tracks slice")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(tracks) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		tracks[i], tracks[j] = tracks[j], tracks[i]
	}

	randomTracks := []string{}
	for _, track := range tracks[:n] {
		randomTracks = append(randomTracks, track.URI)
	}

	return randomTracks, nil
}
