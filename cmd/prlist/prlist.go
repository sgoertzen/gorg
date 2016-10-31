package main

import (
	"github.com/sgoertzen/repoclone"
	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	organization *string
	debug        *bool
}

// Clone all the repos of an orgnaization
func main() {
	c := getConfiguration()
	repoclone.SetDebug(*c.debug)
	repoclone.ListPRs(*c.organization)
}

func getConfiguration() config {
	config := config{}
	config.organization = kingpin.Arg("organization", "GitHub organization to be analyized").Required().String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	kingpin.Version("1.0.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
