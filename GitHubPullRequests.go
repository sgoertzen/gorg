package repoclone

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

// PRSummary is a summary of a PR
type PRSummary struct {
	Repository *string
	Created    *time.Time
	Login      *string
	Title      *string
	Url        *string
}

// ListPullRequests will list all the pull requests for an organization
func ListPullRequests(orgname string) {
	pullRequests := GetPullRequests(orgname)
	printPRsToConsole(pullRequests)
}

// GetPullRequests will return an array of open pull requests
func GetPullRequests(orgname string) *[]PRSummary {
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
		// add to array
		for _, pr := range prs {
			summary := PRSummary{Repository: repo.Name, Created: pr.CreatedAt, Login: pr.User.Login, Title: pr.Title, Url: pr.URL}

			summaries = append(summaries, summary)
		}
	}
	return &summaries
}

func printPRsToConsole(prs *[]PRSummary) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repo", "Created", "Author", "Title", "Link"})

	for _, pr := range *prs {
		table.Append(formatPR(pr))
	}
	table.Render()
}

func formatPR(prSummary PRSummary) []string {
	var formatedTime string

	warnAfter := time.Now().AddDate(0, 0, -3)
	errorAfter := time.Now().AddDate(0, 0, -7)

	format := "2006-01-02"

	created := prSummary.Created
	if created.Before(errorAfter) {
		formatedTime = color.RedString(created.Format(format))
	} else if created.Before(warnAfter) {
		formatedTime = color.YellowString(created.Format(format))
	} else {
		formatedTime = color.GreenString(created.Format(format))
	}

	// TODO: Don't list these.  Find a way to turn a struct into a string array?
	return []string{*prSummary.Repository, formatedTime, *prSummary.Login, *prSummary.Title, *prSummary.Url}
}
