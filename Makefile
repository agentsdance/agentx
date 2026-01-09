VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X github.com/agentsdance/agentx/internal/version.Version=$(VERSION) \
           -X github.com/agentsdance/agentx/internal/version.GitCommit=$(GIT_COMMIT) \
           -X github.com/agentsdance/agentx/internal/version.BuildDate=$(BUILD_DATE)

PREFIX ?= /usr/local

.PHONY: build install clean

build:
	go build -ldflags "$(LDFLAGS)" -o agentx

install: build
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp agentx $(DESTDIR)$(PREFIX)/bin/

clean:
	rm -f agentx
