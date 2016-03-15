// +build integration

package repoclone

import (
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepoNotExistsLocally(t *testing.T) {
	var nonexisting = "doesnotexist"
	repo := github.Repository{FullName: &nonexisting}

	exists := repoExistsLocally(repo)
	assert.False(t, exists)
}

func TestRepoExistsLocally(t *testing.T) {
	var existing = "cmd"
	repo := github.Repository{FullName: &existing}
	exists := repoExistsLocally(repo)
	assert.True(t, exists)
}
