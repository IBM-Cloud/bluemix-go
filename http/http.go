package http

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/IBM-Bluemix/bluemix-go"
)

const (
	userAgentHeader      = "User-Agent"
	authorizationHeader  = "Authorization"
	uaaAccessTokenHeader = "X-Auth-Uaa-Token"

	iamRefreshTokenHeader = "X-Auth-Refresh-Token"
)

//NewHTTPClient ...
func NewHTTPClient(config *bluemix.Config) *http.Client {
	return &http.Client{
		Transport: makeTransport(config),
		Timeout:   config.HTTPTimeout,
	}
}

func makeTransport(config *bluemix.Config) http.RoundTripper {
	return NewTraceLoggingTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   50 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 20 * time.Second,
		DisableCompression:  true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.SSLDisable,
		},
	})
}

//DefaultHeader ...
func DefaultHeader(c *bluemix.Config) http.Header {
	h := http.Header{}
	h.Set(userAgentHeader, UserAgent())
	h.Set(authorizationHeader, c.UAAAccessToken)
	return h
}

//DefaultClusterAuthHeader ...
func DefaultClusterAuthHeader(c *bluemix.Config) http.Header {
	h := http.Header{}
	h.Set(userAgentHeader, UserAgent())
	h.Set(authorizationHeader, c.IAMAccessToken)
	h.Set(iamRefreshTokenHeader, c.IAMRefreshToken)
	h.Set(uaaAccessTokenHeader, c.UAAAccessToken)
	return h
}

//UserAgent ...
func UserAgent() string {
	return fmt.Sprintf("Blumix-go SDK %s / %s ", bluemix.Version, runtime.GOOS)
}
