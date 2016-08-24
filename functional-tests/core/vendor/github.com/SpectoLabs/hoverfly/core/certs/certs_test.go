package certs

import (
	"crypto/x509"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewCert(t *testing.T) {
	x509c, _, err := NewCertificatePair("certy.com", "cert authority", 365*24*time.Hour)
	if err != nil {
		t.Errorf("Failed to generate certificate and key pair, got error: %s", err.Error())
	}

	if err := x509c.VerifyHostname("certy.com"); err != nil {
		t.Errorf("x509c.VerifyHostname(%q): got %v, want no error", "certy.com", err)
	}

	if got, want := x509c.Subject.Organization, []string{"cert authority"}; !reflect.DeepEqual(got, want) {
		t.Errorf("x509c.Subject.Organization: got %v, want %v", got, want)
	}

	if got := x509c.SubjectKeyId; got == nil {
		t.Error("x509c.SubjectKeyId: got nothing, want key ID")
	}
	if !x509c.BasicConstraintsValid {
		t.Error("x509c.BasicConstraintsValid: got false, want true")
	}

	if got, want := x509c.KeyUsage, x509.KeyUsageKeyEncipherment; got&want == 0 {
		t.Error("x509c.KeyUsage: got nothing, want to include x509.KeyUsageKeyEncipherment")
	}
	if got, want := x509c.KeyUsage, x509.KeyUsageDigitalSignature; got&want == 0 {
		t.Error("x509c.KeyUsage: got nothing, want to include x509.KeyUsageDigitalSignature")
	}

	want := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	if got := x509c.ExtKeyUsage; !reflect.DeepEqual(got, want) {
		t.Errorf("x509c.ExtKeyUsage: got %v, want %v", got, want)
	}

	if got, want := x509c.DNSNames, []string{"certy.com"}; !reflect.DeepEqual(got, want) {
		t.Errorf("x509c.DNSNames: got %v, want %v", got, want)
	}
}

func TestNewPriv(t *testing.T) {
	_, priv, err := NewCertificatePair("certy.com", "cert authority", 365*24*time.Hour)
	if err != nil {
		t.Errorf("Failed to generate certificate and key pair, got error: %s", err.Error())
	}

	err = priv.Validate()
	if err != nil {
		t.Errorf("Key validation failed, got error: %s", err.Error())
	}
}

func TestTlsCert(t *testing.T) {
	pub, priv, err := NewCertificatePair("certy.com", "cert authority", 365*24*time.Hour)
	if err != nil {
		t.Errorf("Failed to generate certificate and key pair, got error: %s", err.Error())
	}
	tlsc, err := GetTLSCertificate(pub, priv, "hoverfly.proxy", 365*24*time.Hour)
	if err != nil {
		t.Errorf("Failed to get tls cert, got error: %s", err.Error())
	}

	x509c := tlsc.Leaf
	if x509c == nil {
		t.Fatal("x509c: got nil, want *x509.Certificate")
	}
	if got := x509c.SerialNumber; got.Cmp(MaxSerialNumber) >= 0 {
		t.Errorf("x509c.SerialNumber: got %v, want <= MaxSerialNumber", got)
	}
	if got, want := x509c.Subject.CommonName, "hoverfly.proxy"; got != want {
		t.Errorf("X509c.Subject.CommonName: got %q, want %q", got, want)
	}
	if err := x509c.VerifyHostname("hoverfly.proxy"); err != nil {
		t.Errorf("x509c.VerifyHostname(%q): got %v, want no error", "certy.com", err)
	}

}

func TestGenerateAndSave(t *testing.T) {
	tlsc, err := GenerateAndSave("certy", "cert authority", 1*24*time.Hour)
	if err != nil {
		t.Errorf("Failed to generate tls certificate, got error: %s", err.Error())
	}
	x509c := tlsc.Leaf
	if x509c == nil {
		t.Errorf("x509c: got nil, want *x509.Certificate")
	}

	if _, err := os.Stat("cert.pem"); os.IsNotExist(err) {
		t.Errorf("expected to find it but cert.pem was not created!")
	} else {
		os.Remove("cert.pem")
	}

	if _, err := os.Stat("key.pem"); os.IsNotExist(err) {
		t.Errorf("expected to find it but key.pem was not created!")
	} else {
		os.Remove("key.pem")
	}

}
