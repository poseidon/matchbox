package http

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

// collectVariables collects group selectors, metadata, and request-scoped
// query parameters into a single structured map suitable for rendering
// templates.
func collectVariables(req *http.Request, group *storagepb.Group) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	data["request"] = make(map[string]interface{})
	if group.Metadata != nil {
		err := json.Unmarshal(group.Metadata, &data)
		if err != nil {
			return nil, err
		}
	}
	for key, value := range group.Selector {
		data[strings.ToLower(key)] = value
	}
	// reserved variables
	data["request"] = map[string]interface{}{
		"query":     labelsFromRequest(nil, req),
		"raw_query": req.URL.RawQuery,
	}
	return data, nil
}

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
