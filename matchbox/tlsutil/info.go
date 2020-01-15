package tlsutil

import (
	"crypto/tls"
)

// TLSInfo prepares tls.Config's from TLS filename inputs.
type TLSInfo struct {
	CAFile   string
	CertFile string
	KeyFile  string
}

// ClientConfig returns a tls.Config for client use.
func (info *TLSInfo) ClientConfig() (*tls.Config, error) {
	// CA for verifying the server
	pool, err := NewCertPool([]string{info.CAFile})
	if err != nil {
		return nil, err
	}

	// client certificate (for authentication)
	cert, err := tls.LoadX509KeyPair(info.CertFile, info.KeyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
		// CA bundle the client should trust when verifying a server
		RootCAs: pool,
		// Client certificates to authenticate to the server
		Certificates: []tls.Certificate{cert},
	}, nil
}

// ServerConfig returns a tls.Config for server use.
func (info *TLSInfo) ServerConfig() (*tls.Config, error) {
	// CA for authenticating clients
	pool, err := NewCertPool([]string{info.CAFile})
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		// Call function to load certificate
		GetCertificate: info.GetLoadCertificateFunc(),
		// Client Authentication (required)
		ClientAuth: tls.RequireAndVerifyClientCert,
		// CA for verifying and authorizing client certificates
		ClientCAs: pool,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}, nil
}

// GetLoadCertificateFunc returns a function to load the certificate and private key from disk to allow
// cert/key rotation without requiring the daemon to be restarted
func (info *TLSInfo) GetLoadCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert, err := tls.LoadX509KeyPair(info.CertFile, info.KeyFile)
		if err != nil {
			return nil, err
		}

		return &cert, nil
	}
}
