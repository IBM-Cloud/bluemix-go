package hpcs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIamuumv2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HPCS Suite")
}
