package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/dcarley/dymoscale"
)

type MockScale struct {
	Raw         []byte
	Measurement dymoscale.Measurementer
	Error       error
}

func (m *MockScale) ReadRaw() ([]byte, error) {
	return m.Raw, m.Error
}
func (m *MockScale) ReadMeasurement() (dymoscale.Measurementer, error) {
	return m.Measurement, m.Error
}
func (m *MockScale) Close() error {
	return nil
}

type MockMeasurement struct {
	Result int
	Error  error
}

func (m *MockMeasurement) Grams() (int, error) {
	return m.Result, m.Error
}

var _ = Describe("Dymodump", func() {
	Describe("validateMode", func() {
		It("should return no error for grams mode", func() {
			Expect(validateMode("grams")).To(BeNil())
		})

		It("should return no error for parsed mode", func() {
			Expect(validateMode("parsed")).To(BeNil())
		})

		It("should return no error for raw mode", func() {
			Expect(validateMode("raw")).To(BeNil())
		})

		It("should return error for invalid mode", func() {
			Expect(validateMode("something")).To(MatchError("invalid mode: something"))
		})
	})

	Describe("getResultOrError", func() {
		Context("grams mode", func() {
			var mode = "grams"

			It("should return result as number of grams", func() {
				scale := &MockScale{
					Measurement: &MockMeasurement{
						Result: 23,
					},
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Result: 23"))
			})

			It("should return scale error", func() {
				scale := &MockScale{
					Error: errors.New("something bad"),
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Error: something bad"))
			})

			It("should return measurement error", func() {
				scale := &MockScale{
					Measurement: &MockMeasurement{
						Error: errors.New("something bad"),
					},
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Error: something bad"))
			})
		})

		Context("parsed mode", func() {
			var mode = "parsed"

			It("should return result as struct with field names", func() {
				scale := &MockScale{
					Measurement: &MockMeasurement{
						Result: 23,
					},
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Result: &{Result:23 Error:<nil>}"))
			})

			It("should return scale error", func() {
				scale := &MockScale{
					Error: errors.New("something bad"),
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Error: something bad"))
			})
		})

		Context("raw mode", func() {
			var mode = "raw"

			It("should return result as byte slice", func() {
				scale := &MockScale{
					Raw: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Result: [1 2 3 4 5 6 7 8]"))
			})

			It("should return scale error", func() {
				scale := &MockScale{
					Error: errors.New("something bad"),
				}

				out := getResultOrError(scale, mode)
				Expect(out).To(Equal("Error: something bad"))
			})
		})
	})
})
