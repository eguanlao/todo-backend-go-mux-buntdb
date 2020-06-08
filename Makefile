GOLANGCI_LINT_VERSION=v1.27.0

all: clean lint test build

clean:
	go clean -x

lint:
	# Make sure you have golangci-lint installed first.
	go get -v github.com/golangci/golangci-lint@$(GOLANGCI_LINT_VERSION)
	golangci-lint run -c $(shell go env GOPATH)/pkg/mod/github.com/golangci/golangci-lint\@$(GOLANGCI_LINT_VERSION)/.golangci.yml

test:
	go test -cover ./...

build:
	go build -v