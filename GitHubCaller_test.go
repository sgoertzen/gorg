package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {

	operation := func() error {
		return nil
	}
	start := time.Now()
	makeGitHubCall(operation)
	makeGitHubCall(operation)
	assert.True(t, time.Since(start).Seconds() > 1)
}

func TestErrors(t *testing.T) {

	operation := func() error {
		// nothing
		return errors.New("test")
	}
	err := makeGitHubCall(operation)
	assert.NotNil(t, err)
}
