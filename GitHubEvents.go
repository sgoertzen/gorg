package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

const dateFormat string = "2006-01-02"

// EventSummary represents a repository event, like a PR or branch
type EventSummary struct {
	Repository *string
	LastUsed   *time.Time
	Login      *string
	Title      *string
	URL        *string
}

// EventList is a list of pull requests or branches
type EventList struct {
	summaries *[]EventSummary
}

// GetEvents returns all matching pull requests or branches for the repo
func GetEvents(command string, orgname string, minAge int, maxAge int) *EventList {
	var summaries []EventSummary

	client := getClient()
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if debug {
			log.Printf("Analying repo: %s", *repo.Name)
		}

		switch command {
		case "prs":
			summaries = listPRs(client, repo, orgname)
		case "branches":
			summaries = listBranches(client, repo, orgname)
		}
		if debug {
			log.Printf("Number of Events found: %d", len(summaries))
		}
		summaries = filterEvents(summaries, minAge, maxAge)
	}
	return &EventList{summaries: &summaries}
}
func listPRs(client *github.Client, repo github.Repository, orgname string) []EventSummary {
	var summaries []EventSummary
	opt := &github.PullRequestListOptions{State: "open", Direction: "asc"}
	owner := orgname
	prs, _, err := client.PullRequests.List(owner, *repo.Name, opt)
	check(err)
	for _, pr := range prs {
		summary := EventSummary{Repository: repo.Name, LastUsed: pr.CreatedAt, Login: pr.User.Login, Title: pr.Title, URL: pr.HTMLURL}
		summaries = append(summaries, summary)
	}
	return summaries
}

func listBranches(client *github.Client, repo github.Repository, orgname string) []EventSummary {
	var summaries []EventSummary
	opt := &github.ListOptions{}
	branches, _, err := client.Repositories.ListBranches(orgname, *repo.Name, opt)
	check(err)
	if debug {
		log.Printf("Number of branches found: %d", len(branches))
	}

	for _, branch := range branches {
		commit, _, _ := client.Repositories.GetCommit(orgname, *repo.Name, *branch.Commit.SHA)
		date := commit.Commit.Author.Date
		author := commit.Commit.Author.Name
		url := commit.HTMLURL
		if debug {
			log.Printf("Branch found with name: %s  commit: %s", *branch.Name, *branch.Commit.SHA)
		}
		if *branch.Name == "master" {
			continue
		}
		summary := EventSummary{Repository: repo.Name, Login: author, LastUsed: date, Title: branch.Name, URL: url}
		summaries = append(summaries, summary)
	}
	return summaries
}

func filterEvents(summaries []EventSummary, minAge int, maxAge int) []EventSummary {
	minDate := time.Now().AddDate(0, 0, -minAge)
	maxDate := time.Now().AddDate(0, 0, -maxAge)

	var filteredSummaries []EventSummary
	for _, event := range summaries {
		if event.LastUsed.Before(minDate) && event.LastUsed.After(maxDate) {
			filteredSummaries = append(filteredSummaries, event)
		}
	}
	return filteredSummaries
}

// AsJSON will print the PRList as JSON to the given writer
func (list EventList) AsJSON(writer io.Writer) {
	b, err := json.Marshal(list.summaries)
	check(err)
	writer.Write(b)
}

// AsText will return the projects in a readable test format.
func (list EventList) AsText(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Repo", "Date", "Author", "Title", "Link"})
	for _, e := range *list.summaries {
		table.Append(formatEvent(e, true))
	}
	table.Render()
}

// AsCSV will return the projects in a readable test format.
func (list EventList) AsCSV(w io.Writer) {
	io.WriteString(w, "Repo,Date,Author,Title,Link\n")
	for _, e := range *list.summaries {
		fmt.Fprintf(w, "%s,%s,%s,%s,%s\n", iface(formatEvent(e, false))...)
	}
}

// AsJira will return the projects in a table format used by jira
func (list EventList) AsJira(w io.Writer) {
	io.WriteString(w, "||Repo||Date||Author||Title||Link||\n")
	for _, e := range *list.summaries {
		fmt.Fprintf(w, "|%s|%s|%s|%s|%s|\n", iface(formatEvent(e, false))...)
	}
}

// AsHTML will return the projects in an HTML table format
func (list EventList) AsHTML(w io.Writer) {
	io.WriteString(w, "<html><body><table><thead><tr>")
	io.WriteString(w, "<th>Repo</th><th>Date</th><th>Author</th><th>Title</th><th>Link</th>")
	io.WriteString(w, "</tr></thead><tbody>")
	for _, e := range *list.summaries {
		io.WriteString(w, "<tr>")
		formatted := append(formatEvent(e, false), *e.URL) // Need url in there twice
		fmt.Fprintf(w, "<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td><a href=\"%s\">%s</a></td>", iface(formatted)...)
		io.WriteString(w, "</tr>")
	}
	io.WriteString(w, "</tbody></table></body></html>")
}

func formatEvent(event EventSummary, colorText bool) []string {
	formatedTime := ""

	if event.LastUsed != nil {
		lastUsed := event.LastUsed
		formatedTime = lastUsed.Format(dateFormat)
		if colorText {
			warnAfter := time.Now().AddDate(0, 0, -3)
			errorAfter := time.Now().AddDate(0, 0, -7)

			if lastUsed.Before(errorAfter) {
				formatedTime = color.RedString(formatedTime)
			} else if lastUsed.Before(warnAfter) {
				formatedTime = color.YellowString(formatedTime)
			} else {
				formatedTime = color.GreenString(formatedTime)
			}
		}
	}
	login := ""
	if event.Login != nil {
		login = *event.Login
	}
	title := ""
	if event.Title != nil {
		title = *event.Title
	}
	return []string{*event.Repository, formatedTime, login, title, *event.URL}
}

func iface(list []string) []interface{} {
	vals := make([]interface{}, len(list))
	for i, v := range list {
		vals[i] = v
	}
	return vals
}

func printEvents(prlist *EventList, filename string, format string) {
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w2 := bufio.NewWriter(f)
		print(prlist, w2, format)
		w2.Flush()
		f.Sync()
	} else {
		print(prlist, os.Stdout, format)
	}
}

// TODO: Move this somewhere
func print(prlist *EventList, w io.Writer, format string) {

	switch strings.ToLower(format) {
	case "text":
		prlist.AsText(w)
	case "json":
		prlist.AsJSON(w)
	case "csv":
		prlist.AsCSV(w)
	case "confluence":
		prlist.AsJira(w)
	case "html":
		prlist.AsHTML(w)
	default:
		panic("Unknown format " + format)
	}
}
