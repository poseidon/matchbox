package api

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"golang.org/x/net/context"
)

const ipxeBootstrap = `#!ipxe
chain ipxe?uuid=${uuid}&mac=${net0/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
`

var ipxeTemplate = template.Must(template.New("iPXE config").Parse(`#!ipxe
kernel {{.Kernel}}{{range $key, $value := .Cmdline}} {{if $value}}{{$key}}={{$value}}{{else}}{{$key}}{{end}}{{end}}
initrd {{ range $element := .Initrd }}{{$element}} {{end}}
boot
`))

// ipxeInspect returns a handler that responds with the iPXE script to gather
// client machine data and chainload to the ipxeHandler.
func ipxeInspect() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, ipxeBootstrap)
	}
	return http.HandlerFunc(fn)
}

// ipxeBoot returns a handler which renders the iPXE boot script for the
// requester.
func ipxeHandler() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		spec, err := specFromContext(ctx)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		var buf bytes.Buffer
		err = ipxeTemplate.Execute(&buf, spec.BootConfig)
		if err != nil {
			log.Errorf("error rendering template: %v", err)
			http.NotFound(w, req)
			return
		}
		if _, err := buf.WriteTo(w); err != nil {
			log.Errorf("error writing to response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return ContextHandlerFunc(fn)
}
