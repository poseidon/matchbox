package api

import (
	"net"
	"net/http"
	"strings"
)

// labelsFromRequest returns Labels from request query parameters.
func labelsFromRequest(req *http.Request) Labels {
	values := req.URL.Query()
	labels := map[string]string{}
	for key := range values {
		switch strings.ToLower(key) {
		case "mac":
			// set mac if and only if it parses
			if hw, err := parseMAC(values.Get(key)); err == nil {
				labels[key] = hw.String()
			}
		default:
			// matchers don't use multi-value keys, drop later values
			labels[key] = values.Get(key)
		}
	}
	return LabelSet(labels)
}

// parseMAC wraps net.ParseMAC with logging.
func parseMAC(s string) (net.HardwareAddr, error) {
	macAddr, err := net.ParseMAC(s)
	if err != nil {
		// invalid MAC arguments may be common
		log.Debugf("error parsing MAC address: %s", err)
		return nil, err
	}
	return macAddr, err
}
