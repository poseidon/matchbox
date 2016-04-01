package api

import (
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// pixiecoreHandler returns a handler that renders the boot config JSON for
// the requester, to implement the Pixiecore API specification.
// https://github.com/danderson/pixiecore/blob/master/README.api.md
func pixiecoreHandler(srv server.Server) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		macAddr, err := parseMAC(filepath.Base(req.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// pixiecore only provides MAC addresses
		attrs := map[string]string{"mac": macAddr.String()}
		group, err := srv.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: attrs})
		if err != nil {
			http.NotFound(w, req)
			return
		}
		profile, err := srv.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err != nil {
			http.NotFound(w, req)
			return
		}
		renderJSON(w, profile.Boot)
	}
	return ContextHandlerFunc(fn)
}
