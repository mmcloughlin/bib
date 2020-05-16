package main

import (
	"fmt"
	"net/http"
)

// Links gathers all the URLs from a bibliography.
func Links(b *Bibliography) []string {
	var links []string
	for _, entry := range b.Entries {
		link, ok := entry.Fields["url"]
		if !ok {
			continue
		}
		links = append(links, link.String())
	}
	return links
}

// CheckLink checks whether the given URL exists.
func CheckLink(u string) (err error) {
	r, err := http.Get(u)
	if err != nil {
		return err
	}
	defer func() {
		if errc := r.Body.Close(); errc != nil && err == nil {
			err = errc
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %d", r.StatusCode)
	}

	return nil
}
