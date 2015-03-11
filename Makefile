.PHONY: all godep test build darwin docker

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

docker:
	docker build --force-rm -qt dymoscale .
	docker run --rm -ti \
		--privileged \
		-v /dev/bus/usb:/dev/bus/usb \
		-v ${GOPATH}:/gopath \
		-v ${GOPATH}/bin.linux:/gopath/bin \
		-w /gopath/src/$(shell go list) \
		dymoscale
