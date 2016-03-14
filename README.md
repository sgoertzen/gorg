#RepoClone
Clone or update any number of GitHub repositories in a single command

[![Build Status](https://travis-ci.org/sgoertzen/repoclone.svg?branch=master)](https://travis-ci.org/sgoertzen/repoclone)

## Install:
```
go get github.com/sgoertzen/repoclone/cmd/repoclone
```

## Usage:
```
usage: repoclone [<flags>] <organization>

Flags:
  -?, --help     Show context-sensitive help (also try --help-long and --help-man).
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
