package main

import (
	"os"
	"strings"

	"github.com/sgoertzen/repoclone"
	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	organization *string
	debug        *bool
	format       *string
	minAge       *int
	maxAge       *int
}

func main() {
	c := getConfiguration()
	repoclone.SetDebug(*c.debug)

	prlist := repoclone.GetPullRequests(*c.organization, *c.minAge, *c.maxAge)
	switch strings.ToLower(*c.format) {
	case "text":
		prlist.AsText(os.Stdout)
	case "json":
		prlist.AsJSON(os.Stdout)
	case "csv":
		prlist.AsCSV(os.Stdout)
	case "confluence":
		prlist.AsJira(os.Stdout)
	case "html":
		prlist.AsHTML(os.Stdout)
	default:
		panic("Unknown format " + *c.format)
	}
}

func getConfiguration() config {
	config := config{}
	config.organization = kingpin.Arg("organization", "GitHub organization to be analyized").Required().String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	config.format = kingpin.Flag("format", "Specify the output format.  Should be either 'text', 'json', or 'csv'").Default("text").Short('o').Enum("text", "json", "csv", "confluence", "html")
	config.minAge = kingpin.Flag("minAge", "Show PRs that have been open for this number of days").Default("0").Short('n').Int()
	config.maxAge = kingpin.Flag("maxAge", "Show PRs that have been open less then this number of days").Default("36500").Short('x').Int()
	kingpin.Version("1.0.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
