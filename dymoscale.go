package dymoscale

import (
	"fmt"

	"github.com/kylelemons/gousb/usb"
)

const VendorID usb.ID = 0x0922 // Dymo, all devices

type Scale struct {
	context  *usb.Context
	device   *usb.Device
	endpoint usb.Endpoint
}

// closeWithError closes all outstanding devices and the context, then
// returns the original error.
func closeWithError(context *usb.Context, devices []*usb.Device, err error) (*Scale, error) {
	for _, dev := range devices {
		dev.Close()
	}
	context.Close()

	return nil, err
}

// NewScale opens a connection to a Dymo USB scale. You MUST call Close()
// when you're finished.
func NewScale() (*Scale, error) {
	ctx := usb.NewContext()

	devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if desc.Vendor == VendorID {
			return true
		}

		return false
	})
	if err != nil {
		return closeWithError(ctx, devs, err)
	}

	if len(devs) != 1 {
		err := fmt.Errorf("expected 1 device, found %d", len(devs))
		return closeWithError(ctx, devs, err)
	}

	dev := devs[0]
	ep, err := dev.OpenEndpoint(
		dev.Configs[0].Config,
		dev.Configs[0].Interfaces[0].Number,
		dev.Configs[0].Interfaces[0].Setups[0].Number,
		dev.Configs[0].Interfaces[0].Setups[0].Endpoints[0].Address,
	)
	if err != nil {
		return closeWithError(ctx, devs, err)
	}

	scale := &Scale{
		context:  ctx,
		device:   dev,
		endpoint: ep,
	}

	return scale, nil
}

// ReadRaw gets a raw reading from the scale.
func (s *Scale) ReadRaw() ([]byte, error) {
	buf := make([]byte, s.endpoint.Info().MaxPacketSize)
	_, err := s.endpoint.Read(buf)

	return buf, err
}

// Close closes the USB device and context.
func (s *Scale) Close() {
	s.device.Close()
	s.context.Close()
}
