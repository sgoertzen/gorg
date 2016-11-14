package main

import (
	"github.com/sgoertzen/gorg"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type config struct {
	organization *string
	directory    *string
	debug        *bool
	clone        *bool
	update       *bool
	remove       *bool
}

// Clone all the repos of an orgnaization
func main() {
	c := getConfiguration()
	gorg.SetDebug(*c.debug)
	gorg.Sync(*c.organization, *c.directory, *c.clone, *c.update, *c.remove)
}

func getConfiguration() config {
	wd, _ := os.Getwd()
	config := config{}
	config.organization = kingpin.Arg("organization", "GitHub organization that should be cloned").Required().String()
	config.directory = kingpin.Flag("directory", "Directory where repos are/should be stored").Default(wd).Short('p').String()
	config.debug = kingpin.Flag("debug", "Output debug information during the run.").Default("false").Short('d').Bool()
	config.clone = kingpin.Flag("clone", "Only clone repos (do not update)").Default("true").Short('c').Bool()
	config.update = kingpin.Flag("update", "Only update repos (do not clone).").Default("true").Short('u').Bool()
	config.remove = kingpin.Flag("remove", "Remove local directories that are not in the organization").Default("false").Short('r').Bool()
	kingpin.Version("1.1.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}
