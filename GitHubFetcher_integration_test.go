// +build integration

package repoclone

import (
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepoNotExistsLocally(t *testing.T) {
	var nonexisting = "doesnotexist"
	repo := github.Repository{Name: &nonexisting}

	exists := repoExistsLocally(repo, "./")
	assert.False(t, exists)
}

func TestRepoExistsLocally(t *testing.T) {
	var existing = "cmd"
	repo := github.Repository{Name: &existing}
	exists := repoExistsLocally(repo, "./")
	assert.True(t, exists)
}

// TODO: Need to remove dirs if they are no longer repos in this org
