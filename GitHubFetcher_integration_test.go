// +build integration

package main

import (
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestRepoNotExistsLocally(t *testing.T) {
	var nonexisting = "doesnotexist"
	repo := github.Repository{Name: &nonexisting}

	exists := repoExistsLocally(repo, "./")
	assert.False(t, exists)
}

func TestRepoExistsLocally(t *testing.T) {
	existing := "Existing"
	defer os.RemoveAll("./" + existing)
	os.Mkdir(existing, 0777)

	repo := github.Repository{Name: &existing}
	exists := repoExistsLocally(repo, "./")
	assert.True(t, exists)
}

func TestRemoveExtra(t *testing.T) {
	testDir := "./IntTests/"
	extraDir := testDir + "extraDir/"

	defer os.RemoveAll(testDir)
	os.Mkdir(testDir, 0777)
	os.Mkdir(extraDir, 0777)

	cleanup(testDir, make(map[string]bool))
	assert.False(t, fileExists(extraDir), "Extra directory was not removed")
}

func TestRemoveLeaveExisting(t *testing.T) {
	debug = true
	testDir := "./IntTests/"
	neededDir := testDir + "neededDir/"

	defer os.RemoveAll(neededDir)
	defer os.RemoveAll(testDir)
	os.Mkdir(testDir, 0777)
	os.Mkdir(neededDir, 0777)

	var repos = make(map[string]bool)
	repos["neededDir"] = true
	cleanup(testDir, repos)
	assert.True(t, fileExists(neededDir), "Needed directory was removed")
}
