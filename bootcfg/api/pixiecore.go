package api

import (
	"net/http"
	"path/filepath"

	"github.com/coreos/coreos-baremetal/bootcfg/storage"
)

// pixiecoreHandler returns a handler that renders the boot config JSON for
// the requester, to implement the Pixiecore API specification.
// https://github.com/danderson/pixiecore/blob/master/README.api.md
func pixiecoreHandler(gr *groupsResource, store storage.Store) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		macAddr, err := parseMAC(filepath.Base(req.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// pixiecore only provides MAC addresses
		attrs := map[string]string{"mac": macAddr.String()}
		group, err := gr.findMatch(attrs)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		profile, err := store.ProfileGet(group.Profile)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		renderJSON(w, profile.Boot)
	}
	return http.HandlerFunc(fn)
}
