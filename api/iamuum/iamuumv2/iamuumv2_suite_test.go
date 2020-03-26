package iamuumv2_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIamuumv2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Iamuumv2 Suite")
}
