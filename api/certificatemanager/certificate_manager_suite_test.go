package certificatemanager

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCertificateManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CertificateManager Suite")
}
