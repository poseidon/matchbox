package api

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

const ipxeBootstrap = `#!ipxe
chain ipxe?uuid=${uuid}&mac=${net0/mac:hexhyp}&ip=${ip}&domain=${domain}&hostname=${hostname}&serial=${serial}
`

var ipxeTemplate = template.Must(template.New("ipxe boot").Parse(`#!ipxe
kernel {{.Kernel}}{{range $key, $value := .Cmdline}} {{if $value}}{{$key}}={{$value}}{{else}}{{$key}}{{end}}{{end}}
initrd {{ range $element := .Initrd }}{{$element}} {{end}}
boot
`))

// ipxeInspect returns a handler that responds with an iPXE script to gather
// client machine data and chain load the real boot script.
func ipxeInspect() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, ipxeBootstrap)
	}
	return http.HandlerFunc(fn)
}

// ipxeBoot returns a handler which renders an iPXE boot config script based
// on the machine attribtue query parameters.
func ipxeHandler(store Store) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		attrs := attrsFromRequest(req)
		config, err := store.BootConfig(attrs)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		var buf bytes.Buffer
		err = ipxeTemplate.Execute(&buf, config)
		if err != nil {
			log.Errorf("error rendering template: %v", err)
			http.NotFound(w, req)
			return
		}
		if _, err := buf.WriteTo(w); err != nil {
			log.Errorf("error writing to response: %v", err)
		}
	}
	return http.HandlerFunc(fn)
}
