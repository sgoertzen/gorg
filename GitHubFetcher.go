package repoclone

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var debug bool

// SetDebug turns on debugging output on this library
func SetDebug(d bool) {
	debug = d
}

// CloneRepos clones all repos for an orgnaization
func CloneRepos(orgname string) error {
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if !repoExistsLocally(repo) {
			clone(repo)
		}
	}
	return nil
}

// UpdateRepos updates all repos for an orgnaization
func UpdateRepos(orgname string) error {
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if repoExistsLocally(repo) {
			update(repo)
		}
	}
	return nil
}

// CloneOrUpdateRepos clones or updates all repos for an orgnaization
func CloneOrUpdateRepos(orgname string) error {
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if repoExistsLocally(repo) {
			update(repo)
		} else {
			clone(repo)
		}
	}
	return nil
}

func getAllRepos(orgname string) []github.Repository {
	client := getClient()

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(orgname, opt)
		if err != nil {
			return nil
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	log.Printf("Found %d repo(s) for the organization %s", len(allRepos), orgname)
	return allRepos
}

func getClient() *github.Client {
	var tc *http.Client
	envToken := os.Getenv("GITHUB_TOKEN")
	if len(envToken) > 0 {
		token := oauth2.Token{AccessToken: envToken}
		ts := oauth2.StaticTokenSource(&token)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}
	client := github.NewClient(tc)
	return client
}

func repoExistsLocally(repo github.Repository) bool {
	cwd, _ := os.Getwd()
	fullPath := filepath.Join(cwd, *repo.Name)
	_, err := os.Stat(fullPath)
	return err == nil
}

func update(repo github.Repository) (int, error) {
	if debug {
		log.Printf("Updating %s", *repo.Name)
	}
	directory := "./" + *repo.Name
	return run(directory, "git", "pull")
}

func clone(repo github.Repository) (int, error) {
	if debug {
		log.Printf("Cloning %s (%s)", *repo.Name, *repo.SSHURL)
	}
	return run("./", "git", "clone", *repo.SSHURL)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
