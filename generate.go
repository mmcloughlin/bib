package main

import (
	"io"
	"sort"
	"strings"
	"text/template"
)

//go:generate assets -d templates -o ztemplates.go -map templates

// Generate templated output from the given bibliography and writes to w.
func Generate(w io.Writer, tmpl string, b *Bibliography) error {
	// Parse template.
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}

	// Prepare template data.
	type entry struct {
		Entry

		Formatted string
	}
	type data struct {
		Entries []entry
	}
	d := data{}

	for _, e := range b.Entries {
		f, err := Format(e)
		if err != nil {
			return err
		}
		d.Entries = append(d.Entries, entry{
			Entry:     *e,
			Formatted: f,
		})
	}

	// Execute.
	return t.Execute(w, d)
}

// BuiltinTemplateNames returns names of builtin templates.
func BuiltinTemplateNames() []string {
	names := []string{}
	for k := range templates {
		k = strings.TrimPrefix(k, "/")
		k = strings.TrimSuffix(k, ".tmpl")
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
