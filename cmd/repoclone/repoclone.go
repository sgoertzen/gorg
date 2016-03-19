package main

import (
    "os"
	"github.com/sgoertzen/repoclone"
	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	organization *string
    directory    *string
	debug        *bool
}

// Clone all the repos of an orgnaization
func main() {
    // TODO: Actually use *c.directory
	c := getConfiguration()
    repoclone.SetDebug(*c.debug)
	repoclone.CloneOrUpdateRepos(*c.organization)
}

func getConfiguration() config {
	config := config{}
	config.organization = kingpin.Arg("organization", "GitHub organization that should be cloned").Required().String()
	config.directory = kingpin.Flag("directory", "Directory where repos are/should be stored").Default(os.Getwd()).Short('p').String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	kingpin.Version("1.1.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
