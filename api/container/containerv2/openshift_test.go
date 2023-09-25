package containerv2

/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Schematics
 * (C) Copyright IBM Corp. 2023 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	OauthPrivateEndpoint         = "https://c200-e.private.test.containers.cloud.ibm.com:31579/.well-known/oauth-authorization-server"
	OauthVirtualPrivateEndpoint  = "https://clusterID.vpe.private.test.containers.cloud.ibm.com:31579/.well-known/oauth-authorization-server"
	ServerPrivateEndpoint        = "https://c200.private.test.containers.cloud.ibm.com:31580"
	ServerVirtualPrivateEndpoint = "https://clusterID.vpe.private.test.containers.cloud.ibm.com:31580"
)

var _ = Describe("Openshift utils", func() {
	Describe("Openshift related AuthEndpoint modification", func() {
		Context("AuthEndpoint is configured with private", func() {
			It("should not change it", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: true,
						PrivateServiceEndpointURL:     ServerPrivateEndpoint,
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthPrivateEndpoint, "private", &clusterInfo)
				Expect(authServer).ShouldNot(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
				Expect(authServer).Should(Equal(OauthPrivateEndpoint))
			})
		})
		Context("AuthEndpoint is configured with VPE", func() {
			It("should not change it", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: true,
						PrivateServiceEndpointURL:     ServerPrivateEndpoint,
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthVirtualPrivateEndpoint, VirtualPrivateEndpoint, &clusterInfo)
				Expect(authServer).ShouldNot(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
				Expect(authServer).Should(Equal(OauthVirtualPrivateEndpoint))
			})
		})
		Context("AuthEndpoint is configured with Private", func() {
			It("should replace to VPE - case for ROKS 4.12", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: true,
						PrivateServiceEndpointURL:     ServerPrivateEndpoint,
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthPrivateEndpoint, VirtualPrivateEndpoint, &clusterInfo)
				Expect(authServer).ShouldNot(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
				Expect(authServer).Should(Equal(OauthVirtualPrivateEndpoint))
			})
		})
		Context("AuthEndpoint is configured with VPE", func() {
			It("should replace to Private - case for ROKS 4.13", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: true,
						PrivateServiceEndpointURL:     ServerPrivateEndpoint,
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthVirtualPrivateEndpoint, PrivateServiceEndpoint, &clusterInfo)
				Expect(authServer).ShouldNot(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
				Expect(authServer).Should(Equal(OauthPrivateEndpoint))
			})
		})
		Context("AuthEndpoint is empty", func() {
			It("should shall fail", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: true,
						PrivateServiceEndpointURL:     ServerPrivateEndpoint,
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint("", PrivateEndpointDNS, &clusterInfo)
				Expect(authServer).Should(BeEmpty())
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("AuthEndpoint is configured with private - request to change", func() {
			It("should fail as Virtual Private Endpoint URL is empty", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: false,
						PrivateServiceEndpointURL:     "",
					},
					VirtualPrivateEndpointURL: "",
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthPrivateEndpoint, VirtualPrivateEndpoint, &clusterInfo)
				Expect(authServer).Should(BeEmpty())
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("AuthEndpoint is configured with VPE - request to change", func() {
			It("should fail as Private Service Endpoint URL is empty", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: false,
						PrivateServiceEndpointURL:     "",
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthVirtualPrivateEndpoint, PrivateServiceEndpoint, &clusterInfo)
				Expect(authServer).Should(BeEmpty())
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("AuthEndpoint is configured with VPE", func() {
			It("should not change it as endpointType is empty", func() {
				clusterInfo := ClusterInfo{
					ServiceEndpoints: Endpoints{
						PrivateServiceEndpointEnabled: false,
						PrivateServiceEndpointURL:     "",
					},
					VirtualPrivateEndpointURL: ServerVirtualPrivateEndpoint,
				}
				authServer, err := reconfigureAuthorizationEndpoint(OauthVirtualPrivateEndpoint, "", &clusterInfo)
				Expect(authServer).NotTo(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
				Expect(authServer).Should(Equal(OauthVirtualPrivateEndpoint))
			})
		})
	})
})
