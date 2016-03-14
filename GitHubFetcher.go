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

// CloneAllRepos clones all repos for an orgnaization
func CloneAllRepos(orgname string) error {
	client := getClient()
	// org, resp, _ := client.Organizations.Get(orgname)
	// switch resp.StatusCode {
	// case 200:
	// 	break
	// case 401:
	// 	log.Printf("Access not authorized.  Add your Github token to an environment variable GITHUB_TOKEN")
	// 	return nil
	// case 404:
	// 	log.Printf("No GitHub organization found named %s", orgname)
	// 	return nil
	// default:
	// 	log.Printf("Unknown status from organization lookup: %d", resp.StatusCode)
	// 	return nil
	// }
	// check(github.CheckResponse(resp.Response))
	// log.Printf("Found %d private repo(s) for %s", *org.TotalPrivateRepos, orgname)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(orgname, opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	for _, repo := range allRepos {
		log.Println("Cloning " + *repo.Name + " (" + *repo.SSHURL + ")")
		clone(*repo.SSHURL)
	}
	return nil
}

// RefreshAllRepos gets the latest for all repos
func RefreshAllRepos(org string) error {
	return nil
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

func clone(cloneURL string) (int, error) {
	app, err := exec.LookPath("git")
	cmd := exec.Command(app, "clone", cloneURL)

	path := "./"
	cmd.Dir = path
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
