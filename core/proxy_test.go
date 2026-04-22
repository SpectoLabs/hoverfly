package hoverfly

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	testURL             = "localhost:8888"
	testProxyAuthHeader = "Proxy-Authorization"
	testHost            = "test.com"
	testPath            = "/testing"
	testHostSSL         = "test.com:443"
	testSchemeHTTPS     = "https://"
)

func Test_authFromHeader_ShouldRemoveProxyAuthorizationHeader(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add(testProxyAuthHeader, "something")

	authFromHeader(req, nil, nil)
	Expect(req.Header).ToNot(HaveKey(testProxyAuthHeader))
}

func Test_authFromHeader_ShouldRemoveXHoverflyAuthorizationHeader(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add("X-HOVERFLY-AUTHORIZATION", "something")

	authFromHeader(req, nil, nil)
	Expect(req.Header).ToNot(HaveKey("X-HOVERFLY-AUTHORIZATION"))
}

func Test_authFromHeader_ShouldReturnErrorIfNotBasicOrBearer(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add(testProxyAuthHeader, "Something YmVuamloOlBhc3N3b3JkMTIz")

	err := authFromHeader(req, nil, nil)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("407 Unknown authentication type `Something`, only `Basic` or `Bearer` are supported"))
}

func Test_authFromHeader_Basic_ShouldBase64DecodeUsernameAndPassword(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add(testProxyAuthHeader, "Basic YmVuamloOlBhc3N3b3JkMTIz")

	var basicUsername, basicPassword string

	Expect(authFromHeader(req, func(username, password string) bool {
		basicUsername = username
		basicPassword = password
		return true
	}, nil)).To(BeNil())

	Expect(basicUsername).To(Equal("benjih"))
	Expect(basicPassword).To(Equal("Password123"))
}

func Test_authFromHeader_Basic_ShouldReturnFalseIfNotBase64Encoded(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add(testProxyAuthHeader, "Basic benjih:Password123")

	Expect(authFromHeader(req, nil, nil)).ToNot(BeNil())
}

func Test_authFromHeader_Basic_ShouldReturnFalseIfDecodedBasicCredentialsArentFormattedCorrectly(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add(testProxyAuthHeader, "Basic YmVuamlo")

	Expect(authFromHeader(req, nil, nil)).ToNot(BeNil())
}

func Test_authFromHeader_Bearer_ShouldPassJwtTokenOntoFunction(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Add(testProxyAuthHeader, "Bearer gregg.EEewGREQ.GDSG")

	var bearerToken string

	Expect(authFromHeader(req, nil, func(token string) bool {
		bearerToken = token
		return true
	})).To(BeNil())

	Expect(bearerToken).To(Equal("gregg.EEewGREQ.GDSG"))
}

func Test_NewProxy_ShouldHandleConnectForHttps(t *testing.T) {
	RegisterTestingT(t)
	https := httptest.NewTLSServer(nil)
	testHoverfly := NewHoverfly()
	shouldHandleConnect(t, testHoverfly, https.URL)
}

func Test_NewProxy_ShouldHandleConnectForHttp(t *testing.T) {
	RegisterTestingT(t)
	var httpServer = httptest.NewServer(nil)
	testHoverfly := NewHoverfly()
	testHoverfly.Cfg.PlainHttpTunneling = true
	shouldHandleConnect(t, testHoverfly, httpServer.URL)
}

func shouldHandleConnect(t *testing.T, hoverfly *Hoverfly, url string) {
	proxy := NewProxy(hoverfly)
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()
	conn, err := net.Dial("tcp", proxyServer.Listener.Addr().String())
	if err != nil {
		t.Fatal("dialing to proxy", err)
	}
	connReq, err := http.NewRequest("CONNECT", url, nil)
	if err != nil {
		t.Fatal("create new request", connReq)
	}
	connReq.Write(conn)
	resp, err := http.ReadResponse(bufio.NewReader(conn), connReq)
	Expect(resp.StatusCode).To(Equal(200))
}

func Test_matchesFilter_ShouldMatchHostDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter(testHost)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeTrue())

	httpsResult := matchesFilter(testHost)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchPathDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter(testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeTrue())

	httpsResult := matchesFilter(testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchHostAndPathDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter(testHost+testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeTrue())

	httpsResult := matchesFilter(testHost+testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchSchemeDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter("https")(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeFalse())

	httpsResult := matchesFilter(testSchemeHTTPS)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchSchemeAndHostDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter(testSchemeHTTPS+testHost)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeFalse())

	httpsResult := matchesFilter(testSchemeHTTPS+testHost)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchSchemeAndHostAndPathDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter(testSchemeHTTPS+testHost+testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeFalse())

	httpsResult := matchesFilter(testSchemeHTTPS+testHost+testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHost,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShoulRemoveSslPortFromHostToMatch(t *testing.T) {
	RegisterTestingT(t)

	httpsResult := matchesFilter(testSchemeHTTPS+testHost+testPath)(&http.Request{
		Host: testHostSSL,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHostSSL,
			Path:   testPath,
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())

	noSchemeResult := matchesFilter(testSchemeHTTPS+testHost+testPath)(&http.Request{
		Host: testHostSSL,
		URL: &url.URL{
			Scheme: "",
			Host:   testHostSSL,
			Path:   testPath,
		},
	}, nil)
	Expect(noSchemeResult).To(BeTrue())
}

func Test_matchesFilter_ShoulRemovePathFromFilterIfConnectRequest(t *testing.T) {
	RegisterTestingT(t)

	removedPathResult := matchesFilter(testSchemeHTTPS+testHost+testPath)(&http.Request{
		Host:   testHostSSL,
		Method: http.MethodConnect,
		URL: &url.URL{
			Scheme: "https",
			Host:   testHostSSL,
			Path:   "",
		},
	}, nil)
	Expect(removedPathResult).To(BeTrue())
}

func Test_NewProxy_ShouldReturn400ForRelativeURLWithNoHostHeader(t *testing.T) {
	RegisterTestingT(t)
	testHoverfly := NewHoverfly()
	proxy := NewProxy(testHoverfly)
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	// Use raw TCP to send a relative URL with no Host header — http.DefaultClient
	// always produces an absolute URL, so raw TCP is required to exercise NonproxyHandler.
	conn, err := net.Dial("tcp", proxyServer.Listener.Addr().String())
	Expect(err).To(BeNil())
	defer conn.Close()

	fmt.Fprintf(conn, "GET /some-path HTTP/1.1\r\nConnection: close\r\n\r\n")

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
}

func Test_NewProxy_ShouldHandleRelativeURLWithHostHeader(t *testing.T) {
	RegisterTestingT(t)

	// Upstream target that the proxy will forward to
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer targetServer.Close()

	testHoverfly := NewHoverfly()
	proxy := NewProxy(testHoverfly)
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	targetURL, _ := url.Parse(targetServer.URL)

	// Use raw TCP to send a relative URL with a Host header — this is the HTTP/1.1
	// pattern that was previously returning 500.
	conn, err := net.Dial("tcp", proxyServer.Listener.Addr().String())
	Expect(err).To(BeNil())
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", targetURL.Host)

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	Expect(err).To(BeNil())
	// Must not return the "This is a proxy server" 500 — any non-500 response is acceptable
	Expect(resp.StatusCode).ToNot(Equal(http.StatusInternalServerError))
}

func Test_NewProxy_ShouldReturn500WhenRelativeURLHostIsProxyItself(t *testing.T) {
	RegisterTestingT(t)
	testHoverfly := NewHoverfly()
	proxy := NewProxy(testHoverfly)
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	// Set ProxyPort to match the test server's port so the self-address guard fires.
	_, port, _ := net.SplitHostPort(proxyServer.Listener.Addr().String())
	testHoverfly.Cfg.ProxyPort = port

	// Send a relative URL with Host == the proxy's own address (direct browser-to-proxy hit).
	conn, err := net.Dial("tcp", proxyServer.Listener.Addr().String())
	Expect(err).To(BeNil())
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: localhost:%s\r\nConnection: close\r\n\r\n", port)

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
}

func Test_matchesFilter_ShouldGetHostNameFromRequest(t *testing.T) {
	RegisterTestingT(t)

	httpResult := matchesFilter(testPath)(&http.Request{
		Host: testHost,
		URL: &url.URL{
			Scheme: "http",
			Path:   testPath,
		},
	}, nil)
	Expect(httpResult).To(BeTrue())
}
