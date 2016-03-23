package certs

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"crypto/rand"
	"net"
	"time"
	"bytes"
	"math/big"
	"encoding/pem"
	"crypto/ecdsa"
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
)

// MaxSerialNumber - nothing very original, big number
var MaxSerialNumber = big.NewInt(0).SetBytes(bytes.Repeat([]byte{255}, 20))

// PemBlockForKey - based on key returns a block
func PemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Unable to marshal ECDSA private key")
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

// NewCertificatePair - returns x509 cert + private key
func NewCertificatePair(name, organization string, validity time.Duration) (*x509.Certificate, *rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	pub := priv.Public()

	pkixpub, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}
	h := sha1.New()
	h.Write(pkixpub)
	keyID := h.Sum(nil)

	serial, err := rand.Int(rand.Reader, MaxSerialNumber)
	if err != nil {
		return nil, nil, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{organization},
		},
		SubjectKeyId:          keyID,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now().Add(-validity),
		NotAfter:              time.Now().Add(validity),
		DNSNames:              []string{name},
		IsCA:                  true,
	}

	raw, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, priv)
	if err != nil {
		return nil, nil, err
	}

	x509c, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, nil, err
	}

	return x509c, priv, nil
}

// GetTlsCertificate - takes x509 cert and private key, returns tls.Certificate that is ready for proxy use
func GetTlsCertificate(cert *x509.Certificate, priv *rsa.PrivateKey, hostname string,  validity time.Duration) (*tls.Certificate, error) {
	host, _, err := net.SplitHostPort(hostname)
	if err == nil {
		hostname = host
	}
	pub := priv.Public()

	pkixpub, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	h := sha1.New()
	h.Write(pkixpub)
	keyID := h.Sum(nil)

	serial, err := rand.Int(rand.Reader, MaxSerialNumber)
	if err != nil {
		return nil, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   hostname,
			Organization: cert.Subject.Organization,
		},
		SubjectKeyId:          keyID,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now().Add(validity),
		NotAfter:              time.Now().Add(validity),
	}

	if ip := net.ParseIP(hostname); ip != nil {
		tmpl.IPAddresses = []net.IP{ip}
	} else {
		tmpl.DNSNames = []string{hostname}
	}

	raw, err := x509.CreateCertificate(rand.Reader, tmpl, cert, priv.Public(), priv)
	if err != nil {
		return nil, err
	}

	// Parse certificate bytes to get a leaf certificate
	x509c, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, err
	}

	tlsc := &tls.Certificate{
		Certificate: [][]byte{raw, cert.Raw},
		PrivateKey:  priv,
		Leaf:        x509c,
	}

	return tlsc, nil
}