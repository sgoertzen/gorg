package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsJira(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsJira(buf)
	assert.Contains(t, buf.String(), "|Hello|")
	assert.Contains(t, buf.String(), "|tuser|Hello|http://somewhere.com|")
}

func TestAsCSV(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsCSV(buf)
	assert.Contains(t, buf.String(), "Hello,")
	assert.Contains(t, buf.String(), ",tuser,Hello,http://somewhere.com")
}

func TestAsHTML(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsHTML(buf)
	assert.Contains(t, buf.String(), "<td>tuser</td><td>Hello</td><td><a href=\"http://somewhere.com\">http://somewhere.com</a></td></tr>")
}

func TestAsJSON(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsJSON(buf)
	assert.Contains(t, buf.String(), "\"Login\":\"tuser\",\"Title\":\"Hello\",\"URL\":\"http://somewhere.com\"}]")
}

func TestAsText(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsText(buf)
	assert.Contains(t, buf.String(), "| Hello | ")
	assert.Contains(t, buf.String(), " | tuser  | Hello | http://somewhere.com |")
}

func TestFilter(t *testing.T) {
	list := makeEventListTestData().summaries
	filtered := filterEvents(*list, 0, 100)
	assert.Equal(t, 1, len(filtered), "Event got filtered incorrectly")
}

func TestFilterTooNew(t *testing.T) {
	list := makeEventListTestData().summaries
	filtered := filterEvents(*list, 10, 100)
	assert.Equal(t, 0, len(filtered), "Event got filtered incorrectly")
}

func TestFilterTooOld(t *testing.T) {
	list := makeEventListTestData().summaries
	filtered := filterEvents(*list, 0, 5)
	assert.Equal(t, 0, len(filtered), "Event got filtered incorrectly")
}

func makeEventListTestData() EventList {
	repo := "Hello"
	created := time.Now().AddDate(0, 0, -7)
	login := "tuser"
	url := "http://somewhere.com"
	list := EventList{
		summaries: &[]EventSummary{
			EventSummary{
				Repository: &repo,
				LastUsed:   &created,
				Login:      &login,
				Title:      &repo,
				URL:        &url,
			},
		},
	}
	return list
}
