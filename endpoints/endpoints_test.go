package endpoints

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndPoints", func() {

	Context("When region is us-south", func() {
		locator := newEndpointLocator("us-south")

		It("should return endpoints with region us-south", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.ng.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.ng.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accountmanagement.ng.bluemix.net"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.ng.bluemix.net"))
			Expect(locator.ContainerEndpoint()).To(Equal("https://us-south.containers.bluemix.net"))
		})
	})

	Context("When region is eu-gb", func() {
		locator := newEndpointLocator("eu-gb")

		It("should return endpoints with region eu-gb", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.eu-gb.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.eu-gb.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accountmanagement.eu-gb.bluemix.net"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.eu-gb.bluemix.net"))
		})
	})

	Context("When region is au-syd", func() {
		locator := newEndpointLocator("au-syd")

		It("should return endpoints with region au-syd", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.au-syd.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.au-syd.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accountmanagement.au-syd.bluemix.net"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.au-syd.bluemix.net"))
		})
	})

	Context("When region is eu-de", func() {
		locator := newEndpointLocator("eu-de")

		It("should return endpoints with region eu-de", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.eu-de.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.eu-de.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accountmanagement.eu-de.bluemix.net"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.eu-de.bluemix.net"))
		})
	})

	Context("When region is not supported", func() {
		locator := newEndpointLocator("in")

		It("should return error", func() {
			_, err := locator.CFAPIEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.UAAEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.AccountManagementEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.IAMEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.ContainerEndpoint()
			Expect(err).To(HaveOccurred())

		})
	})

})

func newEndpointLocator(region string) EndpointLocator {
	return NewEndpointLocator(region)
}
