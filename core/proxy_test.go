package hoverfly

import (
	"net/http"
	"net/url"
	"testing"

	"bufio"
	"net"

	. "github.com/onsi/gomega"

	"net/http/httptest"
)

func Test_authFromHeader_ShouldRemoveProxyAuthorizationHeader(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "something")

	authFromHeader(req, nil, nil)
	Expect(req.Header).ToNot(HaveKey("Proxy-Authorization"))
}

func Test_authFromHeader_ShouldRemoveXHoverflyAuthorizationHeader(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("X-HOVERFLY-AUTHORIZATION", "something")

	authFromHeader(req, nil, nil)
	Expect(req.Header).ToNot(HaveKey("X-HOVERFLY-AUTHORIZATION"))
}

func Test_authFromHeader_ShouldReturnErrorIfNotBasicOrBearer(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Something YmVuamloOlBhc3N3b3JkMTIz")

	err := authFromHeader(req, nil, nil)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("407 Unknown authentication type `Something`, only `Basic` or `Bearer` are supported"))
}

func Test_authFromHeader_Basic_ShouldBase64DecodeUsernameAndPassword(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Basic YmVuamloOlBhc3N3b3JkMTIz")

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
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Basic benjih:Password123")

	Expect(authFromHeader(req, nil, nil)).ToNot(BeNil())
}

func Test_authFromHeader_Basic_ShouldReturnFalseIfDecodedBasicCredentialsArentFormattedCorrectly(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Basic YmVuamlo")

	Expect(authFromHeader(req, nil, nil)).ToNot(BeNil())
}

func Test_authFromHeader_Bearer_ShouldPassJwtTokenOntoFunction(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Bearer gregg.EEewGREQ.GDSG")

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
	httpResult := matchesFilter("test.com")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeTrue())

	httpsResult := matchesFilter("test.com")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchPathDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter("/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeTrue())

	httpsResult := matchesFilter("/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchHostAndPathDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter("test.com/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeTrue())

	httpsResult := matchesFilter("test.com/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchSchemeDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter("https")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeFalse())

	httpsResult := matchesFilter("https://")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchSchemeAndHostDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter("https://test.com")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeFalse())

	httpsResult := matchesFilter("https://test.com")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShouldMatchSchemeAndHostAndPathDestination(t *testing.T) {
	RegisterTestingT(t)
	httpResult := matchesFilter("https://test.com/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeFalse())

	httpsResult := matchesFilter("https://test.com/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())
}

func Test_matchesFilter_ShoulRemoveSslPortFromHostToMatch(t *testing.T) {
	RegisterTestingT(t)

	httpsResult := matchesFilter("https://test.com/testing")(&http.Request{
		Host: "test.com:443",
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com:443",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpsResult).To(BeTrue())

	noSchemeResult := matchesFilter("https://test.com/testing")(&http.Request{
		Host: "test.com:443",
		URL: &url.URL{
			Scheme: "",
			Host:   "test.com:443",
			Path:   "/testing",
		},
	}, nil)
	Expect(noSchemeResult).To(BeTrue())
}

func Test_matchesFilter_ShoulRemovePathFromFilterIfConnectRequest(t *testing.T) {
	RegisterTestingT(t)

	removedPathResult := matchesFilter("https://test.com/testing")(&http.Request{
		Host: "test.com:443",
		Method: http.MethodConnect,
		URL: &url.URL{
			Scheme: "https",
			Host:   "test.com:443",
			Path:   "",
		},
	}, nil)
	Expect(removedPathResult).To(BeTrue())
}


func Test_matchesFilter_ShouldGetHostNameFromRequest(t *testing.T) {
	RegisterTestingT(t)

	httpResult := matchesFilter("/testing")(&http.Request{
		Host: "test.com",
		URL: &url.URL{
			Scheme: "http",
			Path:   "/testing",
		},
	}, nil)
	Expect(httpResult).To(BeTrue())
}

