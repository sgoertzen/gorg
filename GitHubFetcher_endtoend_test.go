// +build endtoend

package repoclone

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testDir = "../test_directory/"

func TestMain(m *testing.M) {
    debug = true
    
	run("./cmd/repoclone", "go", "build") // compile repoclone
	run("./", "mkdir", testDir)        // make test dir
	output := m.Run()                     // run tests (in build dir)
	//run("../", "rm", "-rf", testDir)    // remove build dir
	os.Exit(output)
}

func TestUpdate(t *testing.T) {
	//defer os.RemoveAll("fuzzy-octo-parakeet")
	fuzzyDir := testDir + "/fuzzy-octo-parakeet"
    
    
    run(testDir, "../repoclone/cmd/repoclone/repoclone", "RepoFetch")
    
    //run(testDir, "git", "clone", "git@github.com:RepoFetch/fuzzy-octo-parakeet.git")
    run(fuzzyDir, "git", "reset", "840a42c1029c20b7b510753162894f4e47dcde1f")
    run(fuzzyDir, "rm", "SecondFile.txt")
    
	assert.False(t, fileExists("SecondFile.txt"))

	// Do cloneOrUpdate
    run(testDir, "../repoclone/cmd/repoclone/repoclone", "RepoFetch")
	//CloneOrUpdateRepos("RepoFetch", false)

	//assert that SecondFile does exist
	assert.True(t, fileExists(fuzzyDir + "/SecondFile.txt"))

}