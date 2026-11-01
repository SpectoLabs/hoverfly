package hoverfly

import (
	"crypto/tls"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	testClientCert = "../functional-tests/core/testdata/cert.pem"
	testClientKey  = "../functional-tests/core/testdata/key.pem"
)

func Test_GetHttpClient_ClientAuth_VerifiesServerCertWithoutCACert(t *testing.T) {
	RegisterTestingT(t)

	server := newClientAuthTLSServer()
	defer server.Close()

	unit := newClientAuthHoverfly(true, "")

	_, err := getWithClientAuth(unit, server)

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("x509"))
}

func Test_GetHttpClient_ClientAuth_TrustsServerSignedByCACert(t *testing.T) {
	RegisterTestingT(t)

	server := newClientAuthTLSServer()
	defer server.Close()

	unit := newClientAuthHoverfly(true, writeServerCert(t, server))

	body, err := getWithClientAuth(unit, server)

	Expect(err).To(BeNil())
	Expect(body).To(Equal("client cert received"))
}

func Test_GetHttpClient_ClientAuth_RejectsServerNotSignedByCACert(t *testing.T) {
	RegisterTestingT(t)

	server := newClientAuthTLSServer()
	defer server.Close()

	unit := newClientAuthHoverfly(true, testClientCert)

	_, err := getWithClientAuth(unit, server)

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("x509"))
}

func Test_GetHttpClient_ClientAuth_SkipsVerificationWhenTLSVerificationDisabled(t *testing.T) {
	RegisterTestingT(t)

	server := newClientAuthTLSServer()
	defer server.Close()

	for _, caCert := range []string{"", testClientCert} {
		unit := newClientAuthHoverfly(false, caCert)

		body, err := getWithClientAuth(unit, server)

		Expect(err).To(BeNil())
		Expect(body).To(Equal("client cert received"))
	}
}

func Test_GetHttpClient_ReturnsDefaultClientWhenDestinationDoesNotMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := newClientAuthHoverfly(true, "")
	unit.Cfg.ClientAuthenticationDestination = "other.com"

	client, err := GetHttpClient(unit, "test.com")

	Expect(err).To(BeNil())
	Expect(client).To(BeIdenticalTo(unit.HTTP))
}

func newClientAuthTLSServer() *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte("client cert received"))
	}))
	server.TLS = &tls.Config{ClientAuth: tls.RequireAnyClientCert}
	server.StartTLS()
	return server
}

func newClientAuthHoverfly(tlsVerification bool, caCert string) *Hoverfly {
	return NewHoverflyWithConfiguration(&Configuration{
		TLSVerification:                 tlsVerification,
		ClientAuthenticationDestination: ".*",
		ClientAuthenticationClientCert:  testClientCert,
		ClientAuthenticationClientKey:   testClientKey,
		ClientAuthenticationCACert:      caCert,
	})
}

func getWithClientAuth(hf *Hoverfly, server *httptest.Server) (string, error) {
	client, err := GetHttpClient(hf, server.Listener.Addr().String())
	Expect(err).To(BeNil())

	resp, err := client.Get(server.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func writeServerCert(t *testing.T, server *httptest.Server) string {
	certPath := filepath.Join(t.TempDir(), "ca.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	Expect(os.WriteFile(certPath, certPEM, 0600)).To(Succeed())
	return certPath
}
