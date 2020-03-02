
GOLANGCI_LINT_VERSION=v1.23.7

LINTER=./bin/golangci-lint
LINTER_VERSION_FILE=./bin/.golangci-lint-version-$(GOLANGCI_LINT_VERSION)

.PHONY: build clean test lint

build:
	go get ./...
	go build ./...

clean:
	go clean

test: build
	go get -t ./...
	@# Note, we need to specify all these packages individually for go test in order to remain 1.8-compatible
	go test -race -v . ./httphelpers ./ldservices

$(LINTER_VERSION_FILE):
	rm -f $(LINTER)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s $(GOLANGCI_LINT_VERSION)
	touch $(LINTER_VERSION_FILE)

lint: $(LINTER_VERSION_FILE)
	go get ./...
	$(LINTER) run ./...
