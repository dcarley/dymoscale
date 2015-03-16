package dymoscale

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
)

// OuncesToGramsVariation is the amount of precision we lose when converting
// from ounces to grams. This is ±28.35*0.1 grams because ounces are only
// reported to one decimal place.
const OuncesToGramsVariation = 2

var _ = Describe("Dymoscale", func() {
	Describe("Measurement", func() {
		Describe("errors", func() {
			It("should error on short read", func() {
				buf := bytes.NewBuffer([]byte{1, 2})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(MatchError("unexpected EOF"))
				Expect(reading).To(Equal(&Measurement{0, 0, 0, 0, 0, 0}))

				grams, err := reading.Grams()
				Expect(err).To(MatchError(ErrInvalidRead))
				Expect(grams).To(Equal(0))
			})

			It("should error (needs tare) on negative weight", func() {
				buf := bytes.NewBuffer([]byte{3, 5, 11, 255, 0, 0, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   NeedsTare,
					Mode:        Ounces,
					ScaleFactor: -1,
					WeightMinor: 0,
					WeightMajor: 0,
				}))

				grams, err := reading.Grams()
				Expect(err).To(MatchError(ErrNeedsTare))
				Expect(grams).To(Equal(0))
			})
		})

		It("should read zero weight", func() {
			buf := bytes.NewBuffer([]byte{3, 2, 11, 255, 0, 0, 0, 0})

			reading, err := ReadMeasurement(buf)
			Expect(err).To(BeNil())
			Expect(reading).To(Equal(&Measurement{
				AlwaysThree: 3,
				Stability:   NoWeight,
				Mode:        Ounces,
				ScaleFactor: -1,
				WeightMinor: 0,
				WeightMajor: 0,
			}))

			grams, err := reading.Grams()
			Expect(err).To(BeNil())
			Expect(grams).To(Equal(0))
		})

		Describe("136g (4.8oz) object", func() {
			expectedGrams := 136

			It("should read from grams mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 136, 0, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Grams,
					ScaleFactor: 0,
					WeightMinor: 136,
					WeightMajor: 0,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(Equal(expectedGrams))
			})

			It("should read from ounces mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 48, 0, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Ounces,
					ScaleFactor: -1,
					WeightMinor: 48,
					WeightMajor: 0,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(BeNumerically("~", expectedGrams, OuncesToGramsVariation))
			})
		})

		Describe("418g (14.7oz) object", func() {
			expectedGrams := 418

			It("should read from grams mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 162, 1, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Grams,
					ScaleFactor: 0,
					WeightMinor: 162,
					WeightMajor: 1,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(Equal(expectedGrams))
			})

			It("should read from ounces mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 147, 0, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Ounces,
					ScaleFactor: -1,
					WeightMinor: 147,
					WeightMajor: 0,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(BeNumerically("~", expectedGrams, OuncesToGramsVariation))
			})
		})

		Describe("1384g (3lb 0.8oz) object", func() {
			expectedGrams := 1384

			It("should read from grams mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 104, 5, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Grams,
					ScaleFactor: 0,
					WeightMinor: 104,
					WeightMajor: 5,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(Equal(expectedGrams))
			})

			It("should read from ounces mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 232, 1, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Ounces,
					ScaleFactor: -1,
					WeightMinor: 232,
					WeightMajor: 1,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(BeNumerically("~", expectedGrams, OuncesToGramsVariation))
			})
		})

		Describe("2492g (5lb 7.9oz) object", func() {
			expectedGrams := 2492

			It("should read from grams mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 188, 9, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Grams,
					ScaleFactor: 0,
					WeightMinor: 188,
					WeightMajor: 9,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(Equal(expectedGrams))
			})

			It("should read from ounces mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 111, 3, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   Stable,
					Mode:        Ounces,
					ScaleFactor: -1,
					WeightMinor: 111,
					WeightMajor: 3,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(BeNumerically("~", expectedGrams, OuncesToGramsVariation))
			})
		})

		Describe("4294g (9lb 7.5oz) object", func() {
			expectedGrams := 4294

			It("should read from grams mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 2, 0, 198, 16, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   4,
					Mode:        2,
					ScaleFactor: 0,
					WeightMinor: 198,
					WeightMajor: 16,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(Equal(expectedGrams))
			})

			It("should read from ounces mode", func() {
				buf := bytes.NewBuffer([]byte{3, 4, 11, 255, 234, 5, 0, 0})

				reading, err := ReadMeasurement(buf)
				Expect(err).To(BeNil())
				Expect(reading).To(Equal(&Measurement{
					AlwaysThree: 3,
					Stability:   4,
					Mode:        11,
					ScaleFactor: -1,
					WeightMinor: 234,
					WeightMajor: 5,
				}))

				grams, err := reading.Grams()
				Expect(err).To(BeNil())
				Expect(grams).To(BeNumerically("~", expectedGrams, OuncesToGramsVariation))
			})
		})
	})
})
