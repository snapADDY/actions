BINDIR := $(CURDIR)/bin

# go option
TESTFLAGS := -race
LDFLAGS := -w -s

.PHONY: all install-js-dependencies detect-changes prepare-img-tag clean test

all: prepare-img-tag detect-changes


install-js-dependencies:
	pnpm install

detect-changes: install-js-dependencies
	pnpm build --filter ...detect-changes

prepare-img-tag:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/ ./prepare-img-tag

clean:
	@rm -rf '$(BINDIR)'

test:
	go test ./... $(TESTFLAGS)

lint:
	golangci-lint run
