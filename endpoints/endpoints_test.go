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
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accounts.cloud.ibm.com"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.ContainerEndpoint()).To(Equal("https://containers.cloud.ibm.com"))
			Expect(locator.CisEndpoint()).To(Equal("https://api.cis.cloud.ibm.com"))
			Expect(locator.ICDEndpoint()).To(Equal("https://api.us-south.databases.cloud.ibm.com"))
			Expect(locator.GlobalSearchEndpoint()).To(Equal("https://api.global-search-tagging.cloud.ibm.com"))
			Expect(locator.GlobalTaggingEndpoint()).To(Equal("https://tags.global-search-tagging.cloud.ibm.com"))
			Expect(locator.MCCPAPIEndpoint()).To(Equal("https://mccp.us-south.cf.cloud.ibm.com"))
			Expect(locator.IAMPAPEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.ResourceManagementEndpoint()).To(Equal("https://resource-controller.cloud.ibm.com"))
			Expect(locator.ResourceControllerEndpoint()).To(Equal("https://resource-controller.cloud.ibm.com"))
			Expect(locator.ResourceCatalogEndpoint()).To(Equal("https://globalcatalog.cloud.ibm.com"))
			Expect(locator.ContainerRegistryEndpoint()).To(Equal("https://registry.ng.bluemix.net"))
		})
	})

	Context("When region is eu-gb", func() {
		locator := newEndpointLocator("eu-gb")

		It("should return endpoints with region eu-gb", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.eu-gb.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.eu-gb.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accounts.cloud.ibm.com"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.CisEndpoint()).To(Equal("https://api.cis.cloud.ibm.com"))
			Expect(locator.ICDEndpoint()).To(Equal("https://api.eu-gb.databases.cloud.ibm.com"))
			Expect(locator.GlobalSearchEndpoint()).To(Equal("https://api.global-search-tagging.cloud.ibm.com"))
			Expect(locator.GlobalTaggingEndpoint()).To(Equal("https://tags.global-search-tagging.cloud.ibm.com"))
		})
	})

	Context("When region is au-syd", func() {
		locator := newEndpointLocator("au-syd")

		It("should return endpoints with region au-syd", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.au-syd.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.au-syd.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accounts.cloud.ibm.com"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.CisEndpoint()).To(Equal("https://api.cis.cloud.ibm.com"))
			Expect(locator.GlobalTaggingEndpoint()).To(Equal("https://tags.global-search-tagging.cloud.ibm.com"))
			Expect(locator.ICDEndpoint()).To(Equal("https://api.au-syd.databases.cloud.ibm.com"))
		})
	})

	Context("When region is eu-de", func() {
		locator := newEndpointLocator("eu-de")

		It("should return endpoints with region eu-de", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.eu-de.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.eu-de.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accounts.cloud.ibm.com"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.CisEndpoint()).To(Equal("https://api.cis.cloud.ibm.com"))
			Expect(locator.GlobalTaggingEndpoint()).To(Equal("https://tags.global-search-tagging.cloud.ibm.com"))
			Expect(locator.ICDEndpoint()).To(Equal("https://api.eu-de.databases.cloud.ibm.com"))
		})
	})

	Context("When region is us-east", func() {
		locator := newEndpointLocator("us-east")

		It("should return endpoints with region us-east", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.us-east.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.us-east.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accounts.cloud.ibm.com"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.CisEndpoint()).To(Equal("https://api.cis.cloud.ibm.com"))
			Expect(locator.GlobalTaggingEndpoint()).To(Equal("https://tags.global-search-tagging.cloud.ibm.com"))
			Expect(locator.ICDEndpoint()).To(Equal("https://api.us-east.databases.cloud.ibm.com"))
		})
	})

	Context("When region is jp-tok", func() {
		locator := newEndpointLocator("jp-tok")

		It("should return endpoints with region jp-tok", func() {
			Expect(locator.CFAPIEndpoint()).To(Equal("https://api.jp-tok.bluemix.net"))
			Expect(locator.UAAEndpoint()).To(Equal("https://login.jp-tok.bluemix.net/UAALoginServerWAR"))
			Expect(locator.AccountManagementEndpoint()).To(Equal("https://accounts.cloud.ibm.com"))
			Expect(locator.IAMEndpoint()).To(Equal("https://iam.cloud.ibm.com"))
			Expect(locator.CisEndpoint()).To(Equal("https://api.cis.cloud.ibm.com"))
			Expect(locator.GlobalTaggingEndpoint()).To(Equal("https://tags.global-search-tagging.cloud.ibm.com"))
			Expect(locator.ICDEndpoint()).To(Equal("https://api.jp-tok.databases.cloud.ibm.com"))
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
			_, err = locator.MCCPAPIEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.CisEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.GlobalSearchEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.GlobalTaggingEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.IAMPAPEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.ICDEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.ResourceManagementEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.ResourceControllerEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.ResourceCatalogEndpoint()
			Expect(err).To(HaveOccurred())
			_, err = locator.ContainerRegistryEndpoint()
		})
	})

})

func newEndpointLocator(region string) EndpointLocator {
	return NewEndpointLocator(region)
}
