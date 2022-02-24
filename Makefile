BINDIR := $(CURDIR)/bin

# go option
TESTFLAGS := -race
LDFLAGS := -w -s

# Rebuild the binary if any of these files change
.PHONY: all prepare-img-tag clean test

all: prepare-img-tag

prepare-img-tag:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/ ./prepare-img-tag

clean:
	@rm -rf '$(BINDIR)'

test:
	go test ./... $(TESTFLAGS)

lint:
	golangci-lint run
