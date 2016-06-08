package tlsutil

import (
	"crypto/tls"
)

// TLSInfo prepares tls.Config's from TLS file inputs.
type TLSInfo struct {
	CAFile   string
	CertFile string
	KeyFile  string
}

// ClientConfig returns a tls.Config for client use.
func (info *TLSInfo) ClientConfig() (*tls.Config, error) {
	pool, err := NewCertPool([]string{info.CAFile})
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
		// CA bundle the client should trust when verifying a server
		RootCAs: pool,
	}, nil
}

// ServerConfig returns a tls.Config for server use.
func (info *TLSInfo) ServerConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(info.CertFile, info.KeyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		// Certificates the server should present to clients
		Certificates: []tls.Certificate{cert},
	}, nil
}
