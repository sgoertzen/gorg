package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

// PRHistory shows PRs for a single user
type PRHistory struct {
	Login             string
	Count             int
	Comments          int
	Merged            int
	Additions         uint32
	Deletions         uint32
	AdditionsUnmerged uint32
	DeletionsUnmerged uint32
	TotalTimeSeconds  uint64
}

// HistoryList is a list of pull requests or branches
type HistoryList struct {
	histories *[]PRHistory
}

// GetHistory returns stuff
func GetHistory(orgname string, minAge int, maxAge int) *HistoryList {

	m := make(map[string]PRHistory)

	var histories []PRHistory

	client := getClient()
	allRepos := getAllRepos(orgname)

	count := 0
	for _, repo := range allRepos {
		if debug {
			log.Printf("Analying repo: %s", *repo.Name)
		}

		listPRHistory(client, repo, orgname, &m, minAge, maxAge)
		if debug {
			log.Printf("Number of closed PRs found: %d", len(histories))
		}
		count++
		if count > 10 {
			break
		}
	}

	var list []PRHistory
	for _, value := range m {
		list = append(list, value)
	}
	return &HistoryList{histories: &list}
}

//TODO: Combine this with other listPRS and accept filter state
func listPRHistory(client *github.Client, repo github.Repository, orgname string, m *map[string]PRHistory, minAge int, maxAge int) {

	opt := &github.PullRequestListOptions{State: "closed"}
	for {
		prs, resp, err := client.PullRequests.List(orgname, *repo.Name, opt)
		check(err)
		prs = filterPullRequests(prs, minAge, maxAge)
		for _, pr := range prs {
			fullPR, _, err := client.PullRequests.Get(orgname, *repo.Name, *pr.Number)
			check(err)
			author := *pr.User.Login
			history := getOrCreateHistory(author, m)
			populatePRHistory(&history, *fullPR)
			(*m)[author] = history
		}
		if resp.NextPage == 0 {
			break
		}
		if debug {
			log.Printf("Fetching page %d", resp.NextPage)
		}
		opt.ListOptions.Page = resp.NextPage
	}
}
func getOrCreateHistory(author string, m *map[string]PRHistory) PRHistory {

	var history PRHistory
	var exists bool
	history, exists = (*m)[author]
	if !exists {
		history = PRHistory{Count: 0, Comments: 0, Merged: 0, Additions: 0, Deletions: 0, TotalTimeSeconds: 0, AdditionsUnmerged: 0, DeletionsUnmerged: 0}
	}
	return history
}

func populatePRHistory(history *PRHistory, pr github.PullRequest) {

	history.Comments += *pr.Comments
	history.Login = *pr.User.Login
	if *pr.Merged {
		history.Merged++
		history.Additions += uint32(*pr.Additions)
		history.Deletions += uint32(*pr.Deletions)
	} else {
		history.AdditionsUnmerged += uint32(*pr.Additions)
		history.DeletionsUnmerged += uint32(*pr.Deletions)
	}

	history.Count++
	seconds := uint64(pr.ClosedAt.Sub(*pr.CreatedAt).Seconds())
	if seconds < 5184000 { // ignore items over two months
		history.TotalTimeSeconds += seconds
	}
	if debug {
		log.Println(history)
	}
}

func printHistories(historyList *HistoryList, filename string, format string) {
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w2 := bufio.NewWriter(f)
		printHistory(historyList, w2, format)
		w2.Flush()
		f.Sync()
	} else {
		printHistory(historyList, os.Stdout, format)
	}
}

func filterPullRequests(prs []*github.PullRequest, minAge int, maxAge int) []*github.PullRequest {
	minDate := time.Now().AddDate(0, 0, -minAge)
	maxDate := time.Now().AddDate(0, 0, -maxAge)

	var filteredPRs []*github.PullRequest
	for _, pr := range prs {
		if pr.CreatedAt.Before(minDate) && pr.CreatedAt.After(maxDate) {
			filteredPRs = append(filteredPRs, pr)
		}
	}
	return filteredPRs
}

// TODO: Move this somewhere
func printHistory(historyList *HistoryList, w io.Writer, format string) {

	switch strings.ToLower(format) {
	case "text":
		historyList.AsText(w)
	// case "json":
	// 	prlist.AsJSON(w)
	// case "csv":
	// 	prlist.AsCSV(w)
	// case "confluence":
	// 	prlist.AsJira(w)
	// case "html":
	// 	prlist.AsHTML(w)
	default:
		panic("Unknown format " + format)
	}
}

// AsText will return the projects in a readable test format.
func (list HistoryList) AsText(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Login", "Count", "Avg Comments", "Merged", "Avg Addition", "Avg Deletion", "Avg Time", "Additions Abandoned", "Deletions Abandoned"})
	for _, e := range *list.histories {
		table.Append(formatHistory(e, true))
	}
	table.Render()
}

func formatHistory(h PRHistory, colorText bool) []string {
	merged := h.Merged
	unmerged := h.Count - h.Merged

	count := fmt.Sprintf("%d", h.Count)
	comments := fmt.Sprintf("%d", (h.Comments / h.Count))
	mergedPct := fmt.Sprintf("%d%%", (merged * 100 / h.Count))
	avgAdd := fmt.Sprintf("%d", (h.Additions / uint32(merged)))
	avgDel := fmt.Sprintf("%d", (h.Deletions / uint32(merged)))
	avgTime := fmt.Sprintf("%d", (h.TotalTimeSeconds / uint64(h.Count)))
	avgAddUnmerg := fmt.Sprintf("%d", (h.AdditionsUnmerged / uint32(unmerged)))
	avgDelUnmerg := fmt.Sprintf("%d", (h.DeletionsUnmerged / uint32(unmerged)))
	return []string{h.Login, count, comments, mergedPct, avgAdd, avgDel, avgTime, avgAddUnmerg, avgDelUnmerg}
}
