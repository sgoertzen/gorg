// +build integration

package repoclone

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClone(t *testing.T) {
	defer os.RemoveAll("fuzzy-octo-parakeet")

	err := CloneAllRepos("RepoFetch")
	assert.Nil(t, err)

	_, err = os.Stat("fuzzy-octo-parakeet")
	assert.False(t, os.IsNotExist(err))
}
