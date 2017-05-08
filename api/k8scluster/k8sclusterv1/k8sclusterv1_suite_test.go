package k8sclusterv1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestK8sclusterv1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8sclusterv1 Suite")
}
