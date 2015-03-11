.PHONY: all godep test build darwin

all: test build

godep:
	go get github.com/tools/godep
	godep restore

test: godep
	godep go test ./...

build: godep
	cd dymodump && godep go build

darwin: export CGO_CFLAGS = -I/opt/boxen/homebrew/include
darwin: export CGO_LDFLAGS = -L/opt/boxen/homebrew/lib
darwin: all
