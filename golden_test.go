package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestGolden(t *testing.T) {
	ext := ".in"
	pattern := filepath.Join("testdata", "golden", "*"+ext)
	inputs, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range inputs {
		input := input // scopelint
		noext := strings.TrimSuffix(input, ext)
		name := filepath.Base(noext)
		t.Run(name, func(t *testing.T) {
			// Open bibliography.
			b, err := ReadBibliography(noext + ".bib")
			if err != nil {
				t.Fatal(err)
			}

			// Process the file.
			s, err := ParseFile(input)
			if err != nil {
				t.Fatal(err)
			}

			if err := s.Validate(b); err != nil {
				t.Fatal(err)
			}

			got, err := s.Bytes(b)
			if err != nil {
				t.Fatal(err)
			}

			// Update golden file if requested.
			golden := noext + ".golden"
			if *update {
				if err := ioutil.WriteFile(golden, got, 0o666); err != nil {
					t.Fatal(err)
				}
			}

			// Read golden file.
			expect, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}

			// Compare.
			AssertLinesEqual(t, expect, got)
		})
	}
}

func AssertLinesEqual(t *testing.T, expect, got []byte) {
	t.Helper()

	// Break into lines.
	expectlines := Lines(string(expect))
	gotlines := Lines(string(got))

	if len(expectlines) != len(gotlines) {
		t.Fatalf("line number mismatch: got %v expect %v", len(gotlines), len(expectlines))
	}

	for i := range expectlines {
		if expectlines[i] != gotlines[i] {
			t.Errorf("line %d:\n\tgot    = %q\n\texpect = %q", i+1, gotlines[i], expectlines[i])
		}
	}
}

func Lines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}
