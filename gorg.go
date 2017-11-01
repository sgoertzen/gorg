package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	command        *string
	directory      *string
	debug          *bool
	filename       *string
	organization   *string
	clone          *bool
	update         *bool
	remove         *bool
	format         *string
	minAge         *int
	maxAge         *int
	cloneOverHTTPS *bool
}

// Clone all the repos of an orgnaization
func main() {
	c := getConfiguration()
	SetDebug(*c.debug)

	switch *c.command {
	case "clone":
		Sync(*c.organization, *c.directory, *c.clone, *c.update, *c.remove, *c.cloneOverHTTPS)
		//os.Exit(status)
	case "prs", "branches":
		events := GetEvents(*c.command, *c.organization, *c.minAge, *c.maxAge)
		printEvents(events, *c.filename, *c.format)
	case "prhistory":
		histories := GetHistory(*c.organization, *c.minAge, *c.maxAge)
		printHistories(histories, *c.filename, *c.format)
	}
}

func getConfiguration() config {
	wd, _ := os.Getwd()
	config := config{}
	config.command = kingpin.Arg("command", "The command to run on the organization.  Valid options are: clone, prs, branches, prhistory").Required().Enum("clone", "prs", "branches", "prhistory")
	config.organization = kingpin.Arg("organization", "GitHub organization that should be cloned").Required().String()
	config.directory = kingpin.Flag("directory", "Directory where repos are/should be stored").Default(wd).Short('p').String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	config.clone = kingpin.Flag("clone", "Clone repositories not already downloaded.").Default("true").Short('c').Bool()
	config.update = kingpin.Flag("update", "Update repositories that were previously cloned.").Default("true").Short('u').Bool()
	config.remove = kingpin.Flag("remove", "Remove local directories that are not in the organization").Default("false").Short('r').Bool()
	config.filename = kingpin.Flag("filename", "The file in which the output should be stored.  If this is left off the output will be printed to the console").Short('f').String()
	config.format = kingpin.Flag("format", "Specify the output format.  Should be either 'text', 'json', or 'csv'").Default("text").Short('o').Enum("text", "json", "csv", "confluence", "html")
	config.minAge = kingpin.Flag("minAge", "Show PRs that have been open for this number of days").Default("0").Short('n').Int()
	config.maxAge = kingpin.Flag("maxAge", "Show PRs that have been open less then this number of days").Default("36500").Short('x').Int()
	config.cloneOverHTTPS = kingpin.Flag("https", "Clone repositories using HTTPS instead of SSL").Default("false").Short('h').Bool()
	kingpin.Version("2.1.1")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
