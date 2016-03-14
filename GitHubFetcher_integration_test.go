// +build integration

package repoclone

import (
	"github.com/stretchr/testify/assert"
    "os"
	"testing"
)

func TestClone(t *testing.T) {
	err := CloneAllRepos("RepoFetch")
    assert.Nil(t, err)
    
    _, err = os.Stat("fuzzy-octo-parakeet")
    assert.False(t, os.IsNotExist(err))
    
    os.RemoveAll("fuzzy-octo-parakeet")
}
