# dymoscale

[![GoDoc](https://godoc.org/github.com/dcarley/dymoscale?status.svg)](http://godoc.org/github.com/dcarley/dymoscale) [![Build Status](https://travis-ci.org/dcarley/dymoscale.svg?branch=master)](https://travis-ci.org/dcarley/dymoscale)

Go library for reading Dymo USB postal scales.

You can buy them relatively cheap from Ebay.

Tested with a Dymo M5, but will probably work with other models too.

## Developing

### darwin/amd64

To test and build on Darwin.
```
brew install libusb
make darwin
```

You may get an `access denied` error though when using the binary/library
because of [this problem][libusb_osx] with OS X attaching its HID driver.

[libusb_osx]: https://github.com/libusb/libusb/wiki/FAQ#How_can_I_run_libusb_applications_under_Mac_OS_X_if_there_is_already_a_kernel_extension_installed_for_the_device)

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

## Credit

Goes to this article for doing all the hard work of reverse engineering the
binary protocol:

- [Steven T Snyder - Reading a Dymo USB scale using Python](http://steventsnyder.com/reading-a-dymo-usb-scale-using-python/)
