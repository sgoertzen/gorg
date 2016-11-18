// +build endtoend

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	output := m.Run()
	os.Exit(output)
}

func setup() {
	// Set this to true if you want more detail from the tests
	debug = true

	log.Println("Building...")
	run("./", "go", "install") // compile gorg
}

func TestClone(t *testing.T) {
	dir, err := ioutil.TempDir("", "clone")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fuzzyDir := filepath.Join(dir, "fuzzy-octo-parakeet")

	// Run the program to clone the repo
	run(dir, "gorg", "clone", "RepoFetch", "--https")
	assert.True(t, fileExists(filepath.Join(fuzzyDir, "SecondFile.txt")), "Cloned repository files not found.")
}

func TestRemove(t *testing.T) {
	dir, err := ioutil.TempDir("", "testremove")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	invalidRepoPath := filepath.Join(dir, "NotExistingInOrg")
	os.Mkdir(invalidRepoPath, 0777)

	// Run the program to clone the repo, specifying cleanup (-r)
	run(dir, "gorg", "clone", "RepoFetch", "-r", "--https")

	assert.False(t, fileExists(invalidRepoPath))
}

func TestPRListWithDefaults(t *testing.T) {
	dir, err := ioutil.TempDir("", "testpr")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Run the program to clone the repo
	outputFile := "prs_output.txt"
	run(dir, "gorg", "prs", "--filename="+outputFile, "RepoFetch", "-d")
	outputWithPath := filepath.Join(dir, outputFile)
	assert.True(t, fileExists(outputWithPath))
	b, err := ioutil.ReadFile(outputWithPath)
	assert.Nil(t, err)

	expected := "+---------------------+------------+-----------+--------------------------------+---------------------------------------------------------+\n" +
		"|        REPO         |    DATE    |  AUTHOR   |             TITLE              |                          LINK                           |\n" +
		"+---------------------+------------+-----------+--------------------------------+---------------------------------------------------------+\n" +
		"| fuzzy-octo-parakeet | 2016-11-09 | sgoertzen | Sample PR for end to end tests | https://github.com/RepoFetch/fuzzy-octo-parakeet/pull/1 |\n" +
		"|                     |            |           | - DO NOT CLOSE                 |                                                         |\n" +
		"+---------------------+------------+-----------+--------------------------------+---------------------------------------------------------+\n"
	assert.Equal(t, expected, string(b))
}

func TestBranchesWithDefaults(t *testing.T) {
	dir, err := ioutil.TempDir("", "testbranch")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Run the program to clone the repo
	outputFile := "branch_output.txt"
	run(dir, "gorg", "branches", "--filename="+outputFile, "RepoFetch")
	outputWithPath := filepath.Join(dir, outputFile)

	assert.True(t, fileExists(outputWithPath))
	b, err := ioutil.ReadFile(outputWithPath)
	assert.Nil(t, err, "Error when reading in the branch_output file")

	expected := "+---------------------+------------+----------------+----------+--------------------------------------------------------------------------------------------------+\n" +
		"|        REPO         |    DATE    |     AUTHOR     |  TITLE   |                                               LINK                                               |\n" +
		"+---------------------+------------+----------------+----------+--------------------------------------------------------------------------------------------------+\n" +
		"| fuzzy-octo-parakeet | 2016-11-09 | Shawn Goertzen | SamplePR | https://github.com/RepoFetch/fuzzy-octo-parakeet/commit/e8e173dac360ed447801caede05e3c87ee7c8893 |\n" +
		"+---------------------+------------+----------------+----------+--------------------------------------------------------------------------------------------------+\n"
	assert.Equal(t, expected, string(b))
}

func TestUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "update")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fuzzyDir := filepath.Join(dir, "fuzzy-octo-parakeet")
	repoFile := filepath.Join(fuzzyDir, "SecondFile.txt")

	// Run the program to clone the repo
	//ret, err := runWithRetries(dir, "gorg", "clone", "RepoFetch", "--https")
	status, err := run(dir, "git", "clone", "https://github.com/RepoFetch/fuzzy-octo-parakeet.git")
	assert.Equal(t, 0, status, "Non-zero return value from git clone call")
	assert.Nil(t, err, "Unable to perform initial clone")
	assert.True(t, fileExists(repoFile), "File not cloned correctly")

	// Reset the repo to a previous commit
	_, err = run(fuzzyDir, "git", "reset", "840a42c1029c20b7b510753162894f4e47dcde1f")
	assert.Nil(t, err, "Error reseting the repo to a previous version")

	run(fuzzyDir, "rm", repoFile)
	assert.False(t, fileExists(repoFile), "File not deleted when repo reset to earlier commit.")

	// Run the program again to pull this time
	status, err = runWithRetries(dir, "gorg", "clone", "RepoFetch", "-u", "--https", "-d")
	assert.Equal(t, 0, status, "Non-zero return value from gorg call")
	assert.Nil(t, err, "Unable to perform second clone")
	assert.True(t, fileExists(repoFile), "File was not restored during update.")
}
