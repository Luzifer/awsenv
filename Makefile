default:
	go build -ldflags "-X main.version=$(shell git describe --tags || git rev-parse --short HEAD || echo dev)"

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh
