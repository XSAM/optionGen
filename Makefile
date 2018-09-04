default: build

# Go parameters
GOBUILD=go build
GOCLEAN=go clean

# Build info
BUILD_TIME=`date +%FT%T%z`
BUILD_DATE=`date +%F`
GIT_REVISION=`git rev-parse --short HEAD`
GIT_BRANCH=`git rev-parse --symbolic-full-name --abbrev-ref HEAD`
GIT_DIRTY=`git diff-index --quiet HEAD -- || echo "✗-"`
BUILD_INFO=$(BUILD_DATE)-$(GIT_BRANCH)-$(GIT_REVISION)

# Installation path
GOBIN?=${GOPATH}/bin

# LDFLAGS
LDFLAGS=-ldflags "-s -X main.buildTime=$(BUILD_TIME) -X main.gitRevision=$(GIT_DIRTY)$(GIT_REVISION) -X main.gitBranch=$(GIT_BRANCH)"

# Target
.PHONY: build
build:
	CGO_ENABLED=0 $(GOBUILD) $(ARGS) -o bin/optionGen $(LDFLAGS) ./cmd/optionGen/

.PHONY: install
install: build
	cp ./bin/* $(GOBIN)/

.PHONY: uninstall
uninstall:
	rm $(GOBIN)/optionGen

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -r ./bin