package dymoscale

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/dcarley/gousb/usb"
)

const VendorID usb.ID = 0x0922 // Dymo, all devices

var (
	ErrInvalidRead = fmt.Errorf("scale gave invalid reading")
	ErrNeedsTare   = fmt.Errorf("scale reads negative, please tare")
	ErrWrongMode   = fmt.Errorf("scale is in ounces mode, please switch to grams")
)

type Mode int8

const (
	Grams  Mode = 2
	Ounces Mode = 11
)

type Stability int8

const (
	NoWeight  Stability = 2
	NeedsTare Stability = 5
	Stable    Stability = 4
)

type Measurementer interface {
	Grams() (int, error)
}

// Measurement represents a parsed reading from the scale.
type Measurement struct {
	AlwaysThree int8      // Don't know what this is but it's always 3
	Stability   Stability // How accurate the measurement was
	Mode        Mode      // Grams or Ounces
	ScaleFactor int8      // WeightMinor*10^n when Mode is Ounces
	WeightMinor uint8     //
	WeightMajor uint8     // Overflow for WeightMinor, n*256
}

// errors returns an error if the reading is invalid.
func (m *Measurement) errors() error {
	if m.AlwaysThree != 3 || m.Stability == 0 || m.Mode == 0 {
		return ErrInvalidRead
	}

	if m.Stability == NeedsTare {
		return ErrNeedsTare
	}

	if m.Mode != Grams && m.Stability != NoWeight {
		return ErrWrongMode
	}

	return nil
}

// Grams returns the measurement in grams.
func (m *Measurement) Grams() (int, error) {
	if err := m.errors(); err != nil {
		return 0, err
	}

	grams := int(m.WeightMinor) + (int(m.WeightMajor) * 256)
	return grams, nil
}

// ReadMeasurement obtains a Measurement from an io.Reader.
func ReadMeasurement(reader io.Reader) (Measurementer, error) {
	var reading Measurement
	err := binary.Read(reader, binary.LittleEndian, &reading)

	return &reading, err
}

// closeWithError closes all outstanding devices and the context, then
// returns the original error.
func closeWithError(context *usb.Context, devices []*usb.Device, err error) (Scaler, error) {
	for _, dev := range devices {
		dev.Close()
	}
	context.Close()

	return nil, err
}

type Scaler interface {
	ReadRaw() ([]byte, error)
	ReadMeasurement() (Measurementer, error)
	ReadGrams() (int, error)
	Close() error
}

type Scale struct {
	context  *usb.Context
	device   *usb.Device
	endpoint usb.Endpoint
}

// NewScale opens a connection to a Dymo USB scale. You MUST call Close()
// when you're finished.
func NewScale() (Scaler, error) {
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

	// TODO: Rate limit? Log?
	if err == usb.ERROR_PIPE || err == usb.ERROR_TIMEOUT {
		s.device.Reset()
	}

	return buf, err
}

// ReadMeasurement returns a parsed Measurement from the scale.
func (s *Scale) ReadMeasurement() (Measurementer, error) {
	res, err := ReadMeasurement(s.endpoint)

	// TODO: Rate limit? Log?
	if err == usb.ERROR_PIPE || err == usb.ERROR_TIMEOUT {
		s.device.Reset()
	}

	return res, err
}

// ReadGrams returns a reading from the scale in grams.
func (s *Scale) ReadGrams() (int, error) {
	measurement, err := s.ReadMeasurement()
	if err != nil {
		return 0, err
	}

	return measurement.Grams()
}

// Close closes the USB device and context. If there are any errors then the
// inner-most is returned, but both will still attempt to be closed.
func (s *Scale) Close() error {
	errDev := s.device.Close()
	errCtx := s.context.Close()

	if errDev != nil {
		return errDev
	}
	return errCtx
}
