#gorg
Clone or update any number of GitHub repositories in a single command

[![Build Status](https://travis-ci.org/sgoertzen/gorg.svg?branch=master)](https://travis-ci.org/sgoertzen/gorg)

## Install:
```
go get github.com/sgoertzen/gorg/cmd/gorg
```

## Usage:
```
usage: gorg [<flags>] <command> <organization>

The commands are:

    clone     Clone all the repositories
    prs       List all open pull requets
    branches  List all open branches 

gorg clone <organization>

Flags:
  -?, --help               Show context-sensitive help (also try --help-long and --help-man).
  -p, --directory="/Users/sgoertzen/Code/gocode/src/github.com/sgoertzen/gorg"
                           Directory where repos are/should be stored
  -d, --debug              Output debug information during the run.
  -c, --clone              Only clone repos (do not update)
  -u, --update             Only update repos (do not clone).
  -r, --remove             Remove local directories that are not in the organization
  -f, --filename=FILENAME  The file in which the output should be stored. If this is left off the output will be printed to the console
  -o, --format=text        Specify the output format. Should be either 'text', 'json', or 'csv'
  -n, --minAge=0           Show PRs that have been open for this number of days
  -x, --maxAge=36500       Show PRs that have been open less then this number of days
  -v, --version            Show application version.

Args:
  <organization>  GitHub organization that should be cloned
```

#### GitHub Token
If you are accessing private repositories you will need to set an environment variable with your token.  If you don't have a token yet you can get one from here: https://github.com/settings/tokens 
```
export GITHUB_TOKEN='YOUR_TOKEN_HERE'
```

#### Examples
Clone all repos from an organization named "RepoFetch"
```
gorg clone RepoFetch 
```

List open pull requests for an organization named "RepoFetch"
```
gorg prs RepoFetch 
```

List open branches for an organization named "RepoFetch" (coming soon)
```
gorg branches RepoFetch 
```


List open pull requests for an organization named "RepoFetch" within 7 to 10 days output to a file as json
```
gorg prs RepoFetch --maxAge=10 --minAge=7 --filename=text.txt --format=json
```

##Development
### Running integration tests
```
go test -tags=integration
```
### Running end to end tests
```
go test -tags=endtoend
```
Can also run these tests using docker. 
```
./buildDockerImage.sh
docker run -it sgoertzen/gorg /bin/bash -c 'cd /go/src/github.com/sgoertzen/gorg/ && go test -tags=endtoend'
```
Note: These tests will run against GitHub and therefore require an internet connection