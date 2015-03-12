# dymoscale

[![GoDoc](https://godoc.org/github.com/dcarley/dymoscale?status.svg)](http://godoc.org/github.com/dcarley/dymoscale) [![Build Status](https://travis-ci.org/dcarley/dymoscale.svg?branch=master)](https://travis-ci.org/dcarley/dymoscale)

Go library for reading Dymo USB postal scales.

You can buy them relatively cheap from Ebay.

Tested with a Dymo M5, but will probably work with other models too.

## Developing

### darwin/amd64

To test and build on Darwin. Although it seems that libusb reports an
`access denied` error:
```
brew install libusb
make darwin
```

### linux/amd64

To test, build, and run on Linux in Docker:
```
make docker
make
./dymodump/dymodump
```

### linux/arm

To test and cross-compile to ARM, which can run on a Raspberry Pi:
```
make docker
make test arm
```
