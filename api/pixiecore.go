package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

const pixiecorePath = "/v1/boot/"

// pixiecoreHandler returns a handler that renders Boot Configs as JSON to
// implement the Pixiecore API specification.
// https://github.com/danderson/pixiecore/blob/master/README.api.md
func pixiecoreHandler(bootConfigs BootAdapter) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		mac := strings.TrimPrefix(req.URL.String(), pixiecorePath)
		attrs := MachineAttrs{MAC: net.HardwareAddr(mac)}
		log.Infof("pixiecore boot config request for %+v", attrs)
		bootConfig, err := bootConfigs.Get(attrs)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		json.NewEncoder(w).Encode(bootConfig)
	}
	return http.HandlerFunc(fn)
}
