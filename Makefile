VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X github.com/agentsdance/agentx/internal/version.Version=$(VERSION) \
           -X github.com/agentsdance/agentx/internal/version.GitCommit=$(GIT_COMMIT) \
           -X github.com/agentsdance/agentx/internal/version.BuildDate=$(BUILD_DATE)

PREFIX ?= /usr/local
WAILS ?= $(shell go env GOPATH)/bin/wails

.PHONY: build install clean release gui gui-dev

build:
	go build -ldflags "$(LDFLAGS)" -o agentx

install: build
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp agentx $(DESTDIR)$(PREFIX)/bin/

clean:
	rm -f agentx
	rm -rf gui/build/bin

release:
	npm publish --access public

# GUI targets
gui:
	cd gui && $(WAILS) build -ldflags "$(LDFLAGS)"

gui-dev:
	cd gui && $(WAILS) dev

# Cross-platform GUI builds
gui-darwin-amd64:
	cd gui && $(WAILS) build -ldflags "$(LDFLAGS)" -platform darwin/amd64

gui-darwin-arm64:
	cd gui && $(WAILS) build -ldflags "$(LDFLAGS)" -platform darwin/arm64

gui-windows:
	cd gui && $(WAILS) build -ldflags "$(LDFLAGS)" -platform windows/amd64

gui-linux:
	cd gui && $(WAILS) build -ldflags "$(LDFLAGS)" -platform linux/amd64

gui-all: gui-darwin-amd64 gui-darwin-arm64 gui-windows gui-linux
