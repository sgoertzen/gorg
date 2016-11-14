package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/sgoertzen/gorg"
	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	filename     *string
	organization *string
	debug        *bool
	format       *string
	minAge       *int
	maxAge       *int
}

func main() {
	c := getConfiguration()
	gorg.SetDebug(*c.debug)

	prlist := gorg.GetPullRequests(*c.organization, *c.minAge, *c.maxAge)

	if *c.filename != "" {
		f, err := os.Create(*c.filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w2 := bufio.NewWriter(f)
		print(prlist, w2, *c.format)
		w2.Flush()
		f.Sync()
	} else {
		print(prlist, os.Stdout, *c.format)
	}
}

func print(prlist *gorg.PRList, w io.Writer, format string) {

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
	config := config{}
	config.organization = kingpin.Arg("organization", "GitHub organization to be analyized").Required().String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	config.filename = kingpin.Flag("filename", "The file in which the output should be stored.  If this is left off the output will be printed to the console").Short('f').String()
	config.format = kingpin.Flag("format", "Specify the output format.  Should be either 'text', 'json', or 'csv'").Default("text").Short('o').Enum("text", "json", "csv", "confluence", "html")
	config.minAge = kingpin.Flag("minAge", "Show PRs that have been open for this number of days").Default("0").Short('n').Int()
	config.maxAge = kingpin.Flag("maxAge", "Show PRs that have been open less then this number of days").Default("36500").Short('x').Int()
	kingpin.Version("1.0.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
