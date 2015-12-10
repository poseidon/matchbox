package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const pixiecorePath = "/v1/boot/"

// pixiecoreHandler implements the Pixiecore API Server Spec.
func pixiecoreHandler(bootConfigs BootConfigProvider) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		remainder := strings.TrimPrefix(req.URL.String(), pixiecorePath)
		bootConfig, err := bootConfigs.Get(remainder)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(bootConfig)
	}
	return http.HandlerFunc(fn)
}
