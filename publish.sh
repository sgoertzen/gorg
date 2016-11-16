#!/usr/bin/env bash

function finish {
    # Cleanup
    if [ -f "gorg" ]; then
        rm gorg > /dev/null
    fi
    if [ -f "gorg.zip" ]; then
        rm gorg.zip > /dev/null
    fi
}
trap finish EXIT

function printUsage {
    echo "Usage:"
    echo " publish.sh <version>"
}

if [ -z "$1" ] ; then
    echo "No version passed in."
    printUsage;
    exit 1
fi

version=$1

echo "Creating release ${version}"

# Install github-release tool
go get github.com/aktau/github-release

# build our project
go build ./cmd/gorg/gorg.go

# Tag and push to github
git tag ${version} && git push --tags

# Create the release
github-release release -u sgoertzen -r gorg -t ${version}

# Create the binary
zip gorg${version}.zip gorg

# Upload the file
github-release upload -u sgoertzen -r gorg -t ${version} -f gorg${version}.zip -n gorg${version}.zip

echo "Release created!"