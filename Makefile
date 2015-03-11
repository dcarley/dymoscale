.PHONY: all godep test build arm darwin docker

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

arm: export CC = arm-linux-gnueabihf-gcc
arm: export CXX = arm-linux-gnueabihf-g++
arm: export CGO_ENABLED = 1
arm: export GOARCH = arm
arm: export GOARM = 7
arm: build

docker:
	docker build --force-rm -qt dymoscale .
	docker run --rm -ti \
		--privileged \
		-v /dev/bus/usb:/dev/bus/usb \
		-v ${GOPATH}:/gopath \
		-v ${GOPATH}/bin.linux:/gopath/bin \
		-w /gopath/src/$(shell go list) \
		dymoscale
