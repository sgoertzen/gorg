package repoclone

import (
	"bufio"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// CloneRepos clones all repos for an orgnaization
func CloneRepos(orgname string, verbose bool) error {
    allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
        if !repoExistsLocally(repo) {
            clone(repo)
        }
	}
	return nil
}
// UpdateRepos updates all repos for an orgnaization
func UpdateRepos(orgname string, verbose bool) error {
    allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
        if repoExistsLocally(repo) {
            update(repo)
        }
	}
	return nil
}
// CloneOrUpdateRepos clones or updates all repos for an orgnaization
func CloneOrUpdateRepos(orgname string, verbose bool) error {
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
    var fullPath = "./" + *repo.FullName
    log.Printf("Fetching for %s", fullPath)
    _, err := os.Stat(fullPath)
    return err == nil
}

func update(repo github.Repository) (int, error) {
    log.Printf("Updating %s (%s)", *repo.Name , *repo.SSHURL)
	app, err := exec.LookPath("git")
	check(err)
	cmd := exec.Command(app, "pull")
    cmd.Dir = "./" + *repo.FullName
    return runCommand(cmd)
}

func clone(repo github.Repository) (int, error) {
    log.Printf("Cloning %s (%s)", *repo.Name, *repo.SSHURL)
	app, err := exec.LookPath("git")
	check(err)
	cmd := exec.Command(app, "clone", *repo.SSHURL)
	cmd.Dir = "./"
    return runCommand(cmd)
}

func runCommand(cmd *exec.Cmd) (int, error) {
    
	stdout, err := cmd.StdoutPipe()
	check(err)

	err = cmd.Start()
	check(err)

	in := bufio.NewScanner(stdout)

	for in.Scan() {
		// Uncomment if we want to include git output in the logs
		//log.Printf(in.Text())
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
