package repoclone

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

const dateFormat string = "2006-01-02"

// PRSummary is a summary of a PR
type PRSummary struct {
	Repository *string
	Created    *time.Time
	Login      *string
	Title      *string
	URL        *string
}

// PRList is a list of pull requests
type PRList struct {
	summaries *[]PRSummary
}

// GetPullRequests will return an array of open pull requests
func GetPullRequests(orgname string, minAge int, maxAge int) *PRList {
	var summaries []PRSummary

	client := getClient()
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if debug {
			log.Printf("Analying repo: %s", *repo.Name)
		}

		opt := &github.PullRequestListOptions{State: "open", Direction: "asc"}
		owner := orgname
		prs, _, err := client.PullRequests.List(owner, *repo.Name, opt)
		check(err)
		if debug {
			log.Printf("Number of PRs found: %d", len(prs))
		}

		minDate := time.Now().AddDate(0, 0, -minAge)
		maxDate := time.Now().AddDate(0, 0, -maxAge)
		// add to array
		for _, pr := range prs {
			if pr.CreatedAt.Before(minDate) && pr.CreatedAt.After(maxDate) {
				summary := PRSummary{Repository: repo.Name, Created: pr.CreatedAt, Login: pr.User.Login, Title: pr.Title, URL: pr.HTMLURL}
				summaries = append(summaries, summary)
			}
		}
	}
	return &PRList{summaries: &summaries}
}

// AsJSON will print the PRList as JSON to the given writer
func (list PRList) AsJSON(writer io.Writer) {
	b, err := json.Marshal(list.summaries)
	check(err)
	writer.Write(b)
}

// AsText will return the projects in a readable test format.
func (list PRList) AsText(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Repo", "Created", "Author", "Title", "Link"})
	for _, pr := range *list.summaries {
		table.Append(formatPR(pr))
	}
	table.Render()
}

// AsCSV will return the projects in a readable test format.
func (list PRList) AsCSV(w io.Writer) {
	io.WriteString(w, "Repo,Created,Author,Title,Link\n")
	for _, pr := range *list.summaries {
		fmt.Fprintf(w, "%s,%s,%s,%s,%s\n", *pr.Repository, pr.Created.Format(dateFormat), *pr.Login, *pr.Title, *pr.URL)
	}
}

// AsJira will return the projects in a table format used by jira
func (list PRList) AsJira(w io.Writer) {
	io.WriteString(w, "||Repo||Created||Author||Title||Link||\n")
	for _, pr := range *list.summaries {
		fmt.Fprintf(w, "|%s|%s|%s|%s|%s|\n", *pr.Repository, pr.Created.Format(dateFormat), *pr.Login, *pr.Title, *pr.URL)
	}
}

// AsHTML will return the projects in an HTML table format
func (list PRList) AsHTML(w io.Writer) {
	io.WriteString(w, "<html><body><table><thead><tr>")
	io.WriteString(w, "<th>Repo</th><th>Created</th><th>Author</th><th>Title</th><th>Link</th>")
	io.WriteString(w, "</tr></thead><tbody>")
	for _, pr := range *list.summaries {
		io.WriteString(w, "<tr>")
		fmt.Fprintf(w, "<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td><a href=\"%s\">%s</a></td>", *pr.Repository, pr.Created.Format(dateFormat), *pr.Login, *pr.Title, *pr.URL, *pr.URL)
		io.WriteString(w, "</tr>")
	}
	io.WriteString(w, "</tbody></table></body></html>")
}

func formatPR(prSummary PRSummary) []string {
	var formatedTime string

	warnAfter := time.Now().AddDate(0, 0, -3)
	errorAfter := time.Now().AddDate(0, 0, -7)

	created := prSummary.Created
	if created.Before(errorAfter) {
		formatedTime = color.RedString(created.Format(dateFormat))
	} else if created.Before(warnAfter) {
		formatedTime = color.YellowString(created.Format(dateFormat))
	} else {
		formatedTime = color.GreenString(created.Format(dateFormat))
	}
	return []string{*prSummary.Repository, formatedTime, *prSummary.Login, *prSummary.Title, *prSummary.URL}
}
