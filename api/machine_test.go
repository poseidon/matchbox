package api

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttrsFromRequest(t *testing.T) {
	hwAddr, err := net.ParseMAC("52:54:00:89:d8:10")
	assert.Nil(t, err)
	cases := []struct {
		urlString string
		attrs     MachineAttrs
	}{
		{"http://a.io", MachineAttrs{}},
		{"http://a.io?uuid=a1b2c3", MachineAttrs{UUID: "a1b2c3"}},
		{"http://a.io?mac=52:54:00:89:d8:10", MachineAttrs{MAC: hwAddr}},
		{"http://a.io?mac=52-54-00-89-d8-10", MachineAttrs{MAC: hwAddr}},
		{"http://a.io?uuid=a1b2c3&mac=52:54:00:89:d8:10", MachineAttrs{UUID: "a1b2c3", MAC: hwAddr}},
		// leave MAC nil if it does not parse
		{"http://a.io?mac=?:?:?", MachineAttrs{}},
	}
	for _, c := range cases {
		req, err := http.NewRequest("GET", c.urlString, nil)
		assert.Nil(t, err)
		assert.Equal(t, attrsFromRequest(req), c.attrs)
	}
}
