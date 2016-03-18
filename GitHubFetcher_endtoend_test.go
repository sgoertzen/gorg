// +build endtoend

package repoclone

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var testDir = "../test_directory/"

func TestMain(m *testing.M) {
	setup()
	output := m.Run()
	teardown()
	os.Exit(output)
}

func setup() {
	// Set this to true if you want more detail from the tests
	debug = false

	log.Println("Running build on repoclone")
	run("./cmd/repoclone", "go", "build") // compile repoclone
	log.Println("Creating directory")
	run("./", "mkdir", testDir) // make test dir
}

func teardown() {
	run("../", "rm", "-rf", testDir) // remove build dir
}

func TestClone(t *testing.T) {
	fuzzyDir := testDir + "fuzzy-octo-parakeet"
	defer os.RemoveAll(fuzzyDir)

	// Run the program to clone the repo
	run(testDir, "../repoclone/cmd/repoclone/repoclone", "RepoFetch")
	assert.True(t, fileExists(fuzzyDir+"/SecondFile.txt"))
}

func TestUpdate(t *testing.T) {
	defer os.RemoveAll("fuzzy-octo-parakeet")

	fuzzyDir := testDir + "fuzzy-octo-parakeet/"

	// Run the program to clone the repo
	run(testDir, "../repoclone/cmd/repoclone/repoclone", "RepoFetch")

	// Reset the repo to a previous commit
	run(fuzzyDir, "git", "reset", "840a42c1029c20b7b510753162894f4e47dcde1f")
	run(fuzzyDir, "rm", "SecondFile.txt")
	assert.False(t, fileExists("SecondFile.txt"))

	// Run the program again to pull this time
	run(testDir, "../repoclone/cmd/repoclone/repoclone", "RepoFetch", "-d")

	//assert that SecondFile does exist
	assert.True(t, fileExists(fuzzyDir+"/SecondFile.txt"))
}
