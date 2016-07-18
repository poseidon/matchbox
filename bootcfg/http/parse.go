package http

import (
	"net"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
)

// labelsFromRequest returns request query parameters.
func labelsFromRequest(logger *logrus.Logger, req *http.Request) map[string]string {
	values := req.URL.Query()
	labels := map[string]string{}
	for key := range values {
		switch strings.ToLower(key) {
		case "mac":
			// set mac if and only if it parses
			if hw, err := parseMAC(values.Get(key)); err == nil {
				labels[key] = hw.String()
			} else {
				if logger != nil {
					logger.WithFields(logrus.Fields{
						"mac": values.Get(key),
					}).Warningf("ignoring unparseable MAC address: %v", err)
				}
			}
		default:
			// matchers don't use multi-value keys, drop later values
			labels[key] = values.Get(key)
		}
	}
	return labels
}

// parseMAC wraps net.ParseMAC with logging.
func parseMAC(s string) (net.HardwareAddr, error) {
	macAddr, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	return macAddr, err
}
