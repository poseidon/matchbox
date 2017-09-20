package http

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"context"
	"github.com/Sirupsen/logrus"
)

const gpxeBootstrap = `#!gpxe
chain gpxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
`

var gpxeTemplate = template.Must(template.New("gPXE config").Parse(`#!gpxe
kernel {{.Kernel}}{{range $arg := .Args}} {{$arg}}{{end}}
initrd {{ range $element := .Initrd }}{{$element}} {{end}}
boot
`))

// gpxeInspect returns a handler that responds with the gPXE script to gather
// client machine data and chainload to the ipxeHandler.
func gpxeInspect() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, gpxeBootstrap)
	}
	return ContextHandlerFunc(fn)
}

// gpxeBoot returns a handler which renders the gPXE boot script for the
// requester.
func (s *Server) gpxeHandler() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		profile, err := profileFromContext(ctx)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels": labelsFromRequest(nil, req),
			}).Infof("No matching profile")
			http.NotFound(w, req)
			return
		}

		// match was successful
		s.logger.WithFields(logrus.Fields{
			"labels":  labelsFromRequest(nil, req),
			"profile": profile.Id,
		}).Debug("Matched an gPXE config")

		var buf bytes.Buffer
		err = gpxeTemplate.Execute(&buf, profile.Boot)
		if err != nil {
			s.logger.Errorf("error rendering template: %v", err)
			http.NotFound(w, req)
			return
		}
		if _, err := buf.WriteTo(w); err != nil {
			s.logger.Errorf("error writing to response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return ContextHandlerFunc(fn)
}
