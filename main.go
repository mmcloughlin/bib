package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	os.Exit(main1())
}

func main1() int {
	log.SetPrefix("bib: ")
	log.SetFlags(0)
	if err := mainerr(); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}

var (
	bibfile = flag.String("bib", "", "bibliography file")
	write   = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

func mainerr() error {
	flag.Parse()

	b, err := ReadBibliography(*bibfile)
	if err != nil {
		return err
	}

	for _, filename := range flag.Args() {
		if err := process(filename, b); err != nil {
			return err
		}
	}

	return nil
}

func process(filename string, b *Bibliography) error {
	s, err := ParseFile(filename)
	if err != nil {
		return err
	}

	if err := s.Validate(b); err != nil {
		return err
	}

	out, err := s.Bytes(b)
	if err != nil {
		return err
	}

	if *write {
		err = ioutil.WriteFile(filename, out, 0644)
	} else {
		_, err = os.Stdout.Write(out)
	}

	return err
}
