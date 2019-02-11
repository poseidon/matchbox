package http

import (
	"net/http"
	"testing"

	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLabelsFromRequest(t *testing.T) {
	emptyMap := map[string]string{}
	logger, _ := logtest.NewNullLogger()
	cases := []struct {
		urlString string
		labels    map[string]string
	}{
		{"http://a.io", emptyMap},
		{"http://a.io?uuid=a1b2c3", map[string]string{"uuid": "a1b2c3"}},
		{"http://a.io?uuid=a1b2c3", map[string]string{"uuid": "a1b2c3"}},
		{"http://a.io?mac=52:da:00:89:d8:10", map[string]string{"mac": validMACStr}},
		{"http://a.io?mac=52-da-00-89-d8-10", map[string]string{"mac": validMACStr}},
		{"http://a.io?uuid=a1b2c3&mac=52:da:00:89:d8:10", map[string]string{"uuid": "a1b2c3", "mac": validMACStr}},
		// parse and set MAC regardless of query argument case
		{"http://a.io?UUID=a1b2c3&MAC=52:DA:00:89:d8:10", map[string]string{"UUID": "a1b2c3", "MAC": validMACStr}},
		// ignore MAC addresses which do not parse
		{"http://a.io?mac=x:x:x:x:x:x", emptyMap},
	}
	for _, c := range cases {
		req, err := http.NewRequest("GET", c.urlString, nil)
		assert.Nil(t, err)
		assert.Equal(t, c.labels, labelsFromRequest(logger, req))
	}
}
