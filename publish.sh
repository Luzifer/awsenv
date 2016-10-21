#!/bin/sh

VERSION=$(git describe --tags --exact-match)
REPO=$(basename $(pwd))
ARCHS="linux/386 linux/amd64 linux/arm darwin/amd64 darwin/386 windows/amd64 windows/386"

set -e

if [ -z "${VERSION}" ]; then
  echo "No tag present, stopping build now."
  exit 0
fi

if [ -z "${GITHUB_TOKEN}" ]; then
  echo "Please set \$GITHUB_TOKEN environment variable"
  exit 1
fi

set -x

go get github.com/aktau/github-release
go get github.com/mitchellh/gox

github-release release --user Luzifer --repo ${REPO} --tag ${VERSION} --name ${VERSION} --draft || true

gox -ldflags="-X main.version=${VERSION}" -osarch="${ARCHS}"
sha256sum ${REPO}_* > SHA256SUMS

for file in ${REPO}_* SHA256SUMS; do
  github-release upload --user Luzifer --repo ${REPO} --tag ${VERSION} --name ${file} --file ${file}
done
