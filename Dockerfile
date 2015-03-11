FROM dcarley/golang-rpi

RUN apt-get -y update && apt-get -y install libusb-1.0-0-dev libusb-1.0-0-dev:armhf usbutils
