package repoclone

import (
	"bufio"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneRepos clones all repos for an orgnaization
func CloneRepos(orgname string, debug bool) error {
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if !repoExistsLocally(repo) {
			clone(repo, debug)
		}
	}
	return nil
}

// UpdateRepos updates all repos for an orgnaization
func UpdateRepos(orgname string, debug bool) error {
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if repoExistsLocally(repo) {
			update(repo, debug)
		}
	}
	return nil
}

// CloneOrUpdateRepos clones or updates all repos for an orgnaization
func CloneOrUpdateRepos(orgname string, debug bool) error {
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if repoExistsLocally(repo) {
			update(repo, debug)
		} else {
			clone(repo, debug)
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

func update(repo github.Repository, debug bool) (int, error) {
	if debug {
		log.Printf("Updating %s (%s)", *repo.Name, *repo.SSHURL)
	}
	directory := "./" + *repo.FullName
	return runCommand(directory, "git", "pull", *repo.SSHURL, debug)
}

func clone(repo github.Repository, debug bool) (int, error) {
	if debug {
		log.Printf("Cloning %s (%s)", *repo.Name, *repo.SSHURL)
	}
	return runCommand("./", "git", "clone", *repo.SSHURL, debug)
}

// TODO, allow for any number of arguments
func runCommand(directory string, command string, arguement1 string, arguement2 string, debug bool) (int, error) {
	app, err := exec.LookPath(command)
	check(err)
	cmd := exec.Command(app, arguement1, arguement2)
	cmd.Dir = "./"

	stdout, err := cmd.StdoutPipe()
	check(err)
	err = cmd.Start()
	check(err)
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		if debug {
			log.Printf(in.Text())
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
		return 1, err
	}
	return 0, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
