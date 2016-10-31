#RepoClone
Clone or update any number of GitHub repositories in a single command

[![Build Status](https://travis-ci.org/sgoertzen/repoclone.svg?branch=master)](https://travis-ci.org/sgoertzen/repoclone)

## Install:
```
go get github.com/sgoertzen/repoclone/cmd/repoclone
go install github.com/sgoertzen/repoclone/cmd/repoclone
```

## Usage:
```
usage: repoclone [<flags>] <organization>

Flags:
  -?, --help     Show context-sensitive help (also try --help-long and --help-man).
  -p, --directory="/Code/gocode/src/github.com/sgoertzen/repoclone"
                 Directory where repos are/should be stored
  -d, --debug    Output debug information during the run.
  -c, --clone    Only clone repos (do not update)
  -u, --update   Only update repos (do not clone).
  -r, --remove   Remove local directories that are not in the organization
  -v, --version  Show application version.

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
repoclone RepoFetch 
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
Note: These tests will run against GitHub and therefore require an internet connection