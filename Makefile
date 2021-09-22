PKG    = github.com/benjojo/ssh-key-confirmer
PREFIX = /usr

all: build/ssh-key-confirmer

# NOTE: This repo uses Go modules, and uses a synthetic GOPATH at
# $(CURDIR)/.gopath that is only used for the build cache. $GOPATH/src/ is
# empty.
GO            = GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS =
GO_LDFLAGS    = -s -w

build/ssh-key-confirmer: *.go
	$(GO) install $(GO_BUILDFLAGS) -ldflags "$(GO_LDFLAGS)" .

install: build/ssh-key-confirmer
	install -D -m 0755 build/ssh-key-confirmer "$(DESTDIR)$(PREFIX)/bin/ssh-key-confirmer"

vendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: install vendor