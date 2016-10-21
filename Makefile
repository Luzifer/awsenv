default:
	go build -ldflags "-X main.version=$(shell git describe --tags || git rev-parse --short HEAD || echo dev)"

ci:
	./publish.sh
