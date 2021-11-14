# Go parameters
	GOCMD=go
	GOBUILD=$(GOCMD) build
	GOCLEAN=$(GOCMD) clean
	GOTEST=$(GOCMD) test

all: test

.PHONY: test
test:
	$(GOTEST) -v `go list ./...` -coverprofile=coverage.out
	# go tool cover -html=coverage.out

