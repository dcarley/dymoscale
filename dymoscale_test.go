package dymoscale

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
)

var _ = Describe("Dymoscale", func() {
	Describe("Measurement", func() {
		It("should error on short read", func() {
			buf := bytes.NewBuffer([]byte{1, 2})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(MatchError("unexpected EOF"))
			Expect(reading).To(Equal(Measurement{0, 0, 0, 0, 0, 0}))
		})

		It("should read zero weight", func() {
			buf := bytes.NewBuffer([]byte{3, 2, 11, 255, 0, 0, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   2,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 0,
				WeightMajor: 0,
			}))
		})

		It("should read negative weight (needs tare)", func() {
			buf := bytes.NewBuffer([]byte{3, 5, 11, 255, 0, 0, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   5,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 0,
				WeightMajor: 0,
			}))
		})

		It("should read 136g object", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 136, 0, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        2,
				ScaleFactor: 0,
				WeightMinor: 136,
				WeightMajor: 0,
			}))
		})

		It("should read 4.8oz object (136g)", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 48, 0, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 48,
				WeightMajor: 0,
			}))
		})

		It("should read 418g object", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 162, 1, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        2,
				ScaleFactor: 0,
				WeightMinor: 162,
				WeightMajor: 1,
			}))
		})

		It("should read 14.7oz object (418g)", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 147, 0, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 147,
				WeightMajor: 0,
			}))
		})

		It("should read 1384g object", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 104, 5, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        2,
				ScaleFactor: 0,
				WeightMinor: 104,
				WeightMajor: 5,
			}))
		})

		It("should read 3lb 0.8oz object (1384g)", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 232, 1, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 232,
				WeightMajor: 1,
			}))
		})

		It("should read 2492g object", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 188, 9, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        2,
				ScaleFactor: 0,
				WeightMinor: 188,
				WeightMajor: 9,
			}))
		})

		It("should read 5lb 7.9oz object (2492g)", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 111, 3, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 111,
				WeightMajor: 3,
			}))
		})

		It("should read 4294g object", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 198, 16, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        2,
				ScaleFactor: 0,
				WeightMinor: 198,
				WeightMajor: 16,
			}))
		})

		It("should read 9lb 7.5oz object (4294g)", func() {
			buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 234, 5, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(Measurement{
				AlwaysThree: 3,
				Stability:   4,
				Mode:        11,
				ScaleFactor: -1,
				WeightMinor: 234,
				WeightMajor: 5,
			}))
		})
	})
})
