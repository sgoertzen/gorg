package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	command      *string
	directory    *string
	debug        *bool
	filename     *string
	organization *string
	clone        *bool
	update       *bool
	remove       *bool
	format       *string
	minAge       *int
	maxAge       *int
}

// Clone all the repos of an orgnaization
func main() {
	c := getConfiguration()
	SetDebug(*c.debug)

	switch *c.command {
	case "clone":
		Sync(*c.organization, *c.directory, *c.clone, *c.update, *c.remove)
	case "prs":
		prlist := GetPullRequests(*c.organization, *c.minAge, *c.maxAge)
		printPRs(prlist, *c.filename, *c.format)
	case "branches":
		// TODO output branches
	}

}

// TODO: move this somewhere
func printPRs(prlist *PRList, filename string, format string) {
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w2 := bufio.NewWriter(f)
		print(prlist, w2, format)
		w2.Flush()
		f.Sync()
	} else {
		print(prlist, os.Stdout, format)
	}
}

// TODO: Move this somewhere
func print(prlist *PRList, w io.Writer, format string) {

	switch strings.ToLower(format) {
	case "text":
		prlist.AsText(w)
	case "json":
		prlist.AsJSON(w)
	case "csv":
		prlist.AsCSV(w)
	case "confluence":
		prlist.AsJira(w)
	case "html":
		prlist.AsHTML(w)
	default:
		panic("Unknown format " + format)
	}
}

func getConfiguration() config {
	wd, _ := os.Getwd()
	config := config{}
	config.command = kingpin.Arg("command", "The command to run on the organization").Required().Enum("clone", "prs", "branches")
	config.organization = kingpin.Arg("organization", "GitHub organization that should be cloned").Required().String()
	config.directory = kingpin.Flag("directory", "Directory where repos are/should be stored").Default(wd).Short('p').String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	config.clone = kingpin.Flag("clone", "Only clone repos (do not update)").Default("true").Short('c').Bool()
	config.update = kingpin.Flag("update", "Only update repos (do not clone).").Default("true").Short('u').Bool()
	config.remove = kingpin.Flag("remove", "Remove local directories that are not in the organization").Default("false").Short('r').Bool()
	config.filename = kingpin.Flag("filename", "The file in which the output should be stored.  If this is left off the output will be printed to the console").Short('f').String()
	config.format = kingpin.Flag("format", "Specify the output format.  Should be either 'text', 'json', or 'csv'").Default("text").Short('o').Enum("text", "json", "csv", "confluence", "html")
	config.minAge = kingpin.Flag("minAge", "Show PRs that have been open for this number of days").Default("0").Short('n').Int()
	config.maxAge = kingpin.Flag("maxAge", "Show PRs that have been open less then this number of days").Default("36500").Short('x').Int()
	kingpin.Version("2.0.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
