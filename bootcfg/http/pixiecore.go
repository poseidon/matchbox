package http

import (
	"net/http"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/mikeynap/coreos-baremetal/bootcfg/server"
	pb "github.com/mikeynap/coreos-baremetal/bootcfg/server/serverpb"
)

// pixiecoreHandler returns a handler that renders the boot config JSON for
// the requester, to implement the Pixiecore API specification.
// https://github.com/danderson/pixiecore/blob/master/README.api.md
func (s *Server) pixiecoreHandler(core server.Server) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		// pixiecore only provides a MAC address label
		macAddr, err := parseMAC(filepath.Base(req.URL.Path))
		if err != nil {
			s.logger.Errorf("unparseable MAC address: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		attrs := map[string]string{"mac": macAddr.String()}

		group, err := core.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: attrs})
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"label": macAddr,
			}).Infof("No matching group")
			http.NotFound(w, req)
			return
		}

		profile, err := core.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"label": macAddr,
				"group": group.Id,
			}).Infof("No profile named: %s", group.Profile)
			http.NotFound(w, req)
			return
		}

		// match was successful
		s.logger.WithFields(logrus.Fields{
			"label":   macAddr,
			"group":   group.Id,
			"profile": profile.Id,
		}).Debug("Matched a Pixiecore config")

		s.renderJSON(w, profile.Boot)
	}
	return ContextHandlerFunc(fn)
}
