package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/subcommands"
)

func main() {
	os.Exit(main1())
}

func main1() int {
	base := command{
		Log: log.New(os.Stderr, "bib: ", 0),
	}
	subcommands.Register(&process{command: base}, "")
	subcommands.Register(&linkcheck{command: base}, "")
	subcommands.Register(subcommands.HelpCommand(), "")

	flag.Parse()
	ctx := context.Background()
	return int(subcommands.Execute(ctx))
}

// command is a base for subcommands, providing some basic conveniences.
type command struct {
	Log *log.Logger
}

// UsageError logs a usage error and returns a suitable exit code.
func (c command) UsageError(format string, args ...interface{}) subcommands.ExitStatus {
	c.Log.Printf(format, args...)
	return subcommands.ExitUsageError
}

// Fail logs an error message and returns a failing exit code.
func (c command) Fail(format string, args ...interface{}) subcommands.ExitStatus {
	c.Log.Printf(format, args...)
	return subcommands.ExitFailure
}

// Error logs err and returns a failing exit code.
func (c command) Error(err error) subcommands.ExitStatus {
	return c.Fail(err.Error())
}

// process subcommand.
type process struct {
	command

	bibfile string
	write   bool
}

func (*process) Name() string     { return "process" }
func (*process) Synopsis() string { return "generate bibliography comments" }
func (*process) Usage() string {
	return `Usage: bib process [-w] -bib <bibfile> <source> ...

Generate references comments for citations in given source files.

`
}

func (cmd *process) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.bibfile, "bib", "", "bibliography file")
	f.BoolVar(&cmd.write, "w", false, "write result to (source) files instead of stdout")
}

func (cmd *process) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if cmd.bibfile == "" {
		return cmd.UsageError("must provide bibliography file")
	}

	b, err := ReadBibliography(cmd.bibfile)
	if err != nil {
		return cmd.Error(err)
	}

	for _, filename := range f.Args() {
		if err := cmd.file(filename, b); err != nil {
			return cmd.Error(err)
		}
	}

	return subcommands.ExitSuccess
}

// file processes a single file.
func (cmd *process) file(filename string, b *Bibliography) error {
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

	if cmd.write {
		err = ioutil.WriteFile(filename, out, 0644)
	} else {
		_, err = os.Stdout.Write(out)
	}

	return err
}

// linkcheck subcommand.
type linkcheck struct {
	command

	bibfile string
	verbose bool
}

func (*linkcheck) Name() string     { return "linkcheck" }
func (*linkcheck) Synopsis() string { return "check whether all urls exist" }
func (*linkcheck) Usage() string {
	return `Usage: bib linkcheck [-v] -bib <bibfile>

Check whether all URLs in the database exist.

`
}

func (cmd *linkcheck) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.bibfile, "bib", "", "bibliography file")
	f.BoolVar(&cmd.verbose, "v", false, "verbose output")
}

func (cmd *linkcheck) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if cmd.bibfile == "" {
		return cmd.UsageError("must provide bibliography file")
	}

	b, err := ReadBibliography(cmd.bibfile)
	if err != nil {
		return cmd.Error(err)
	}

	// Check all URLs.
	status := subcommands.ExitSuccess
	for _, link := range Links(b) {
		if err := CheckLink(link); err != nil {
			cmd.Log.Printf("error: %s: %s", link, err)
			status = subcommands.ExitFailure
		} else if cmd.verbose {
			cmd.Log.Printf("ok: %s", link)
		}
	}

	return status
}
