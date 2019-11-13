package schematics_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestK8schematics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8schematics Suite")
}
