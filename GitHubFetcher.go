package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var debug bool

// SetDebug turns on debugging output on this library
func SetDebug(d bool) {
	debug = d
}

// Sync pull down all repos for an orgnaization
func Sync(orgname string, path string, clone bool, update bool, remove bool, useHTTPS bool) {
	allRepos := getAllRepos(orgname)
	allReposMap := make(map[string]bool)
	for _, repo := range allRepos {
		if !repoExistsLocally(repo, path) {
			if clone {
				_, err := doClone(repo, path, useHTTPS)
				check(err)
			}
		} else if update {
			_, err := doUpdate(repo, path)
			check(err)
		}
		allReposMap[*repo.Name] = true
	}
	if remove {
		cleanup(path, allReposMap)
	}
}

func cleanup(path string, repos map[string]bool) {
	files := getDirectories(path)

	for _, directory := range files {
		if !directory.IsDir() {
			if debug {
				log.Printf("Skipping %s as it is not a directory", directory.Name())
			}
			continue
		}
		// check if the directory exists in the map
		if _, b := repos[directory.Name()]; !b {
			if debug {
				log.Printf("Removing %s", directory.Name())
			}
			os.RemoveAll(filepath.Join(path, directory.Name()))
		} else if debug {
			log.Printf("Skipping %s as it is found in the organization", directory.Name())
		}
	}
}

func getDirectories(path string) []os.FileInfo {

	// Get a list of directories off this
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Error reading the directory: " + path)
	}
	return files
}

func getAllRepos(orgname string) []github.Repository {
	client := getClient()

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []github.Repository
	retries := 0
	for {
		repos, resp, err := client.Repositories.ListByOrg(orgname, opt)
		if err != nil {
			if debug {
				log.Printf("Retrying list repos (attempt %d)", retries)
			}
			if retries < 5 {
				time.Sleep(waitTime)
				retries++
				continue
			}
			log.Printf("Error while fetching repos: %s", err)
			return nil
		}
		for _, repo := range repos {
			allRepos = append(allRepos, *repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	if debug {
		log.Printf("Found %d repo(s) for the organization %s", len(allRepos), orgname)
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

func repoExistsLocally(repo github.Repository, path string) bool {
	fullPath := filepath.Join(path, *repo.Name)
	_, err := os.Stat(fullPath)
	return err == nil
}

func doUpdate(repo github.Repository, path string) (int, error) {
	if debug {
		log.Printf("Updating %s", *repo.Name)
	}
	directory := filepath.Join(path, *repo.Name)
	return runWithRetries(directory, "git", "pull")
}

func doClone(repo github.Repository, path string, useHTTPS bool) (int, error) {
	var url string
	if useHTTPS {
		url = *repo.CloneURL
	} else {
		url = *repo.SSHURL
	}
	if debug {
		log.Printf("Cloning %s (%s)", *repo.Name, url)
	}
	return runWithRetries(path, "git", "clone", url)
}

func check(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}
