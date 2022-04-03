PACKAGES=$(shell go list ./... | grep -v 'tests')
HERUMI= $(shell pwd)/.herumi
CGO_LDFLAGS=CGO_LDFLAGS="-L$(HERUMI)/bls/lib -lbls384_256 -lm -lstdc++ -g -O2"
BUILD_LDFLAGS= -ldflags "-X github.com/zarbchain/zarb-wallet/version.build=`git rev-parse --short=8 HEAD`"
RELEASE_LDFLAGS= -ldflags "-s -w"

all: install test


########################################
### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45
	go install github.com/gordonklaus/ineffassign

herumi:
	@if [ ! -d $(HERUMI) ]; then \
		git clone --recursive https://github.com/herumi/bls.git $(HERUMI)/bls && cd $(HERUMI)/bls && make minimized_static; \
	fi


########################################
### Building
build:
	go build $(BUILD_LDFLAGS) -o ./build/zarb-wallet ./cmd

install:
	go install $(BUILD_LDFLAGS) ./cmd

release: herumi
	$(CGO_LDFLAGS) go build $(RELEASE_LDFLAGS) ./cmd

########################################
### Testing
test:
	go test ./... -covermode=atomic

########################################
### Formatting, linting, and vetting
fmt:
	gofmt -s -w .
	golangci-lint run -e "SA1019" \
		--timeout=5m0s \
		--enable=gofmt \
		--enable=unconvert \
		--enable=unparam \
		--enable=revive \
		--enable=asciicheck \
		--enable=misspell \
		--enable=gosec
	ineffassign ./...

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install release
.PHONY: devtools test herumi fmt
