// +build endtoend

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDir = "./test_directory/"
var outputFile = "output.txt"

func TestMain(m *testing.M) {
	setup()
	output := m.Run()
	debug = true
	teardown()
	os.Exit(output)
}

func setup() {
	// Set this to true if you want more detail from the tests
	debug = true

	log.Println("Building...")
	run("./", "mkdir", testDir)             // make test dir
	run("./", "go", "build", testDir+"../") // compile gorg
}

func teardown() {
	run("./", "rm", "-rf", testDir) // remove build dir
}

func TestClone(t *testing.T) {
	fuzzyDir := testDir + "fuzzy-octo-parakeet"
	defer os.RemoveAll(fuzzyDir)

	// Run the program to clone the repo
	run(testDir, "gorg", "clone", "RepoFetch")
	assert.True(t, fileExists(fuzzyDir+"/SecondFile.txt"))
}

func TestPath(t *testing.T) {
	fuzzyDir := testDir + "fuzzy-octo-parakeet"
	defer os.RemoveAll(fuzzyDir)

	// Run the program to clone the repo
	Sync("RepoFetch", testDir, true, true, false)

	run(testDir, "gorg", "clone", "RepoFetch")
	assert.True(t, fileExists(fuzzyDir+"/SecondFile.txt"))
}

func TestUpdate(t *testing.T) {
	defer os.RemoveAll("fuzzy-octo-parakeet")

	fuzzyDir := testDir + "fuzzy-octo-parakeet/"

	// Run the program to clone the repo
	run(testDir, "gorg", "clone", "RepoFetch")

	// Reset the repo to a previous commit
	run(fuzzyDir, "git", "reset", "840a42c1029c20b7b510753162894f4e47dcde1f")
	run(fuzzyDir, "rm", "SecondFile.txt")
	assert.False(t, fileExists("SecondFile.txt"))

	// Run the program again to pull this time
	run(testDir, "gorg", "clone", "RepoFetch", "-d")

	//assert that SecondFile does exist
	assert.True(t, fileExists(fuzzyDir+"/SecondFile.txt"))
}

func TestRemove(t *testing.T) {
	defer os.RemoveAll("fuzzyo-octoo-parakeet")
	invalidRepoPath := testDir + "NotExistingInOrg"
	os.Mkdir(invalidRepoPath, 0777)

	// Run the program to clone the repo, specifying cleanup (-r)
	run(testDir, "gorg", "clone", "RepoFetch", "-r")

	assert.False(t, fileExists(invalidRepoPath))
}

func TestPRListWithDefaults(t *testing.T) {
	// Run the program to clone the repo
	run(testDir, "gorg", "prs", "--filename="+outputFile, "RepoFetch")
	outputWithPath := testDir + outputFile
	assert.True(t, fileExists(outputWithPath))
	b, err := ioutil.ReadFile(outputWithPath)
	assert.Nil(t, err)

	actual := "+---------------------+------------+-----------+--------------------------------+---------------------------------------------------------+\n" +
		"|        REPO         |  CREATED   |  AUTHOR   |             TITLE              |                          LINK                           |\n" +
		"+---------------------+------------+-----------+--------------------------------+---------------------------------------------------------+\n" +
		"| fuzzy-octo-parakeet | 2016-11-09 | sgoertzen | Sample PR for end to end tests | https://github.com/RepoFetch/fuzzy-octo-parakeet/pull/1 |\n" +
		"|                     |            |           | - DO NOT CLOSE                 |                                                         |\n" +
		"+---------------------+------------+-----------+--------------------------------+---------------------------------------------------------+\n"
	assert.Equal(t, actual, string(b))
}
