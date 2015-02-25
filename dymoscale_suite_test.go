package dymoscale

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDymoscale(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dymoscale Suite")
}
