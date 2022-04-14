package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/nickng/bibtex"
)

// Entry in a bibliography.
type Entry struct {
	bibtex.BibEntry
}

// Authors returns the list of authors.
func (e Entry) Authors() []string {
	field, found := e.Fields["author"]
	if !found {
		return nil
	}
	authors := strings.Split(field.String(), " and ")
	for i := range authors {
		authors[i] = strings.TrimSpace(authors[i])
	}
	return authors
}

// DateField parses a field as a date in ISO 8601 format.
func (e Entry) DateField(name string) (time.Time, error) {
	s, ok := e.Fields[name]
	if !ok {
		return time.Time{}, errors.New("field not found")
	}
	return time.Parse("2006-01-02", s.String())
}

// ByCiteName sorts a list of entries by their citation name.
type ByCiteName []*Entry

func (e ByCiteName) Len() int           { return len(e) }
func (e ByCiteName) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ByCiteName) Less(i, j int) bool { return e[i].CiteName < e[j].CiteName }

// Bibliography is a collection of references.
type Bibliography struct {
	Entries []*Entry
}

// ReadBibliography reads entries from the given BiBTeX file.
func ReadBibliography(path string) (b *Bibliography, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if errc := f.Close(); err == nil && errc != nil {
			err = errc
		}
	}()

	bib, err := bibtex.Parse(f)
	if err != nil {
		return nil, err
	}

	// Build.
	b = &Bibliography{}
	for _, e := range bib.Entries {
		if err := b.AddEntry(&Entry{BibEntry: *e}); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// AddEntry adds an entry to the bibliography.
func (b *Bibliography) AddEntry(e *Entry) error {
	if b.Lookup(e.CiteName) != nil {
		return fmt.Errorf("key %q already in bibliography", e.CiteName)
	}
	b.Entries = append(b.Entries, e)
	return nil
}

// Lookup reference with the given key.
func (b *Bibliography) Lookup(key string) *Entry {
	for _, e := range b.Entries {
		if e.CiteName == key {
			return e
		}
	}
	return nil
}

// FormatBibTeX outputs b in a canonical format.
func FormatBibTeX(b *Bibliography) []byte {
	// Sort entries.
	sorted := make([]*Entry, len(b.Entries))
	copy(sorted, b.Entries)
	sort.Sort(ByCiteName(sorted))

	// Convert to bibtex package type and use pretty printing.
	bib := bibtex.NewBibTex()
	for _, entry := range sorted {
		bib.AddEntry(&entry.BibEntry)
	}

	return []byte(bib.PrettyString())
}
