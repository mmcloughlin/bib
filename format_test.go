package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nickng/bibtex"
)

type TestEntry struct {
	Type   string
	Name   string
	Fields map[string]string
}

func (t TestEntry) Entry() *Entry {
	e := bibtex.NewBibEntry(t.Type, t.Name)
	for name, value := range t.Fields {
		e.AddField(name, bibtex.BibConst(value))
	}
	return &Entry{*e}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		TestEntry
		Expect string
	}{
		{
			TestEntry: TestEntry{
				Name: "single_author",
				Type: "misc",
				Fields: map[string]string{
					"author": "First Author",
					"title":  "Title",
				},
			},
			Expect: "First Author. Title.",
		},
		{
			TestEntry: TestEntry{
				Name: "two_authors",
				Type: "misc",
				Fields: map[string]string{
					"author": "First Author and Second Author",
					"title":  "Title",
				},
			},
			Expect: "First Author and Second Author. Title.",
		},
		{
			TestEntry: TestEntry{
				Name: "multi_author",
				Type: "misc",
				Fields: map[string]string{
					"author": "First Author and Second Author and Third Author",
					"title":  "Title",
				},
			},
			Expect: "First Author, Second Author and Third Author. Title.",
		},
		{
			TestEntry: TestEntry{
				Name: "url_urldate",
				Type: "misc",
				Fields: map[string]string{
					"author":  "First Author",
					"title":   "Title",
					"url":     "https://golang.org",
					"urldate": "2020-02-06",
				},
			},
			Expect: "First Author. Title. https://golang.org (accessed February 6, 2020)",
		},
		{
			TestEntry: TestEntry{
				Name: "misc",
				Type: "misc",
				Fields: map[string]string{
					"author":       "First Author",
					"title":        "Title",
					"howpublished": "School Newsletter",
					"license":      "MIT License",
				},
			},
			Expect: "First Author. Title. School Newsletter. MIT License.",
		},
		{
			TestEntry: TestEntry{
				Name: "inproceedings",
				Type: "inproceedings",
				Fields: map[string]string{
					"author":    "First Author",
					"title":     "Title",
					"booktitle": "Handbook of Golang",
					"pages":     "42--78",
				},
			},
			Expect: "First Author. Title. In Handbook of Golang, pages 42--78.",
		},
		{
			TestEntry: TestEntry{
				Name: "article",
				Type: "article",
				Fields: map[string]string{
					"author":  "First Author",
					"title":   "Title",
					"journal": "Journal of Computer Science",
					"year":    "1997",
				},
			},
			Expect: "First Author. Title. Journal of Computer Science. 1997.",
		},
		{
			TestEntry: TestEntry{
				Name: "inbook",
				Type: "inbook",
				Fields: map[string]string{
					"author":    "First Author",
					"title":     "Title",
					"booktitle": "Book Title",
					"chapter":   "7",
					"year":      "1988",
				},
			},
			Expect: "First Author. Title. Book Title, chapter 7. 1988.",
		},
		{
			TestEntry: TestEntry{
				Name: "phdthesis",
				Type: "phdthesis",
				Fields: map[string]string{
					"author": "First Author",
					"title":  "Title",
					"school": "University of Bristol",
					"year":   "2001",
				},
			},
			Expect: "First Author. Title. PhD thesis, University of Bristol. 2001.",
		},
		{
			TestEntry: TestEntry{
				Name: "mastersthesis",
				Type: "mastersthesis",
				Fields: map[string]string{
					"author": "First Author",
					"title":  "Title",
					"school": "University of Bristol",
					"year":   "2001",
				},
			},
			Expect: "First Author. Title. Masters thesis, University of Bristol. 2001.",
		},
		{
			TestEntry: TestEntry{
				Name: "techreport",
				Type: "techreport",
				Fields: map[string]string{
					"author":      "First Author",
					"title":       "Title",
					"institution": "Initech",
					"number":      "TPS-142",
					"year":        "1973",
				},
			},
			Expect: "First Author. Title. Technical Report TPS-142, Initech. 1973.",
		},
	}
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			e := c.Entry()
			got, err := Format(e)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.Expect {
				t.Logf("entry  = %s", e)
				t.Logf("got    = %s", got)
				t.Logf("expect = %s", c.Expect)
				t.Fail()
			}
		})
	}
}

func TestFormatRequiredFields(t *testing.T) {
	cases := []struct {
		Type     string
		Required []string
	}{
		{
			Type:     "misc",
			Required: []string{"title"},
		},
		{
			Type:     "inproceedings",
			Required: []string{"title", "booktitle"},
		},
		{
			Type:     "article",
			Required: []string{"title", "journal"},
		},
		{
			Type:     "inbook",
			Required: []string{"title", "booktitle", "chapter"},
		},
		{
			Type:     "phdthesis",
			Required: []string{"title", "school"},
		},
		{
			Type:     "mastersthesis",
			Required: []string{"title", "school"},
		},
		{
			Type:     "techreport",
			Required: []string{"title", "institution", "number"},
		},
	}
	for _, c := range cases {
		c := c // scopelint
		for _, required := range c.Required {
			required := required // scopelint
			t.Run(c.Type+"_"+required, func(t *testing.T) {
				// Build an entry with all required fields except this one.
				e := TestEntry{
					Name:   required,
					Type:   c.Type,
					Fields: map[string]string{},
				}
				for _, name := range c.Required {
					if name != required {
						e.Fields[name] = name
					}
				}

				// Confirm it errors.
				_, err := Format(e.Entry())
				if err == nil {
					t.Fatal("expected error; got nil")
				}
				substring := fmt.Sprintf("missing required field %q", required)
				if !strings.Contains(err.Error(), substring) {
					t.Logf("error              = %s", err)
					t.Logf("expected substring = %q", substring)
					t.Fail()
				}
			})
		}
	}
}

func TestFormatUnknownType(t *testing.T) {
	e := TestEntry{
		Name:   "unknowntype",
		Type:   "unknowntype",
		Fields: map[string]string{"title": "Title"},
	}
	_, err := Format(e.Entry())
	expect := `unknown entry type "unknowntype"`
	if err == nil || err.Error() != expect {
		t.Fatalf("got error %v; expected %s", err, expect)
	}
}

func TestWrap(t *testing.T) {
	cases := []struct {
		Text   string
		Width  int
		Expect []string
	}{
		{
			Text:  "the fat cat sat on the mat",
			Width: 10,
			Expect: []string{
				"the fat",
				"cat sat on",
				"the mat",
			},
		},
		{
			Text:  "thisismuchlongerthanwidth cat on mat",
			Width: 10,
			Expect: []string{
				"thisismuchlongerthanwidth",
				"cat on mat",
			},
		},
	}
	for _, c := range cases {
		if got := Wrap(c.Text, c.Width); !reflect.DeepEqual(got, c.Expect) {
			t.Errorf("Wrap(%q, %v) = %#v; expect %#v", c.Text, c.Width, got, c.Expect)
		}
	}
}
