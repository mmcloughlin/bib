package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testenv"
	"github.com/rogpeppe/go-internal/testscript"
)

var net = flag.Bool("net", false, "allow network calls")

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"bib": main1,
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("testdata", "scripts"),
		Condition: func(cond string) (bool, error) {
			switch cond {
			case "network":
				return *net && testenv.HasExternalNetwork(), nil
			default:
				return false, fmt.Errorf("unknown condition %q", cond)
			}
		},
	})
}
