package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoryAsText(t *testing.T) {
	list := makeHistoryTestData()
	buf := new(bytes.Buffer)
	list.AsText(buf)

	results := buf.String()
	assert.Contains(t, results, "| tuser ")
	assert.Contains(t, results, " 4 |")
	assert.Contains(t, results, "| 75%  ")
	assert.Contains(t, results, " 5 |")
	assert.Contains(t, results, " 342 |")
	assert.Contains(t, results, " 112 |")
	assert.Contains(t, results, " 3821 |")
	assert.Contains(t, results, " 1234 |")
	assert.Contains(t, results, " 5678 |")
}

func makeHistoryTestData() HistoryList {
	history := HistoryList{
		histories: &[]PRHistory{
			PRHistory{
				Login:             "tuser",
				Count:             4,
				Comments:          21,
				Merged:            3,
				Additions:         1026,
				Deletions:         337,
				TotalTimeSeconds:  15284,
				AdditionsUnmerged: 1234,
				DeletionsUnmerged: 5678,
			},
		},
	}
	return history
}
