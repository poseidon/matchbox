package api

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

const ipxeBootstrap = `#!ipxe
chain config?uuid=${uuid}
`

var ipxeTemplate = template.Must(template.New("ipxe boot").Parse(`#!ipxe
kernel {{.Kernel}} cloud-config-url=cloud/config?uuid=${uuid} {{range $key, $value := .Cmdline}} {{if $value}}{{$key}}={{$value}}{{else}}{{$key}}{{end}}{{end}}
initrd {{ range $element := .Initrd }} {{$element}}{{end}}
boot
`))

// ipxeMux handles iPXE requests for boot (config) scripts.
func ipxeMux(bootConfigs BootAdapter) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/ipxe/boot.ipxe", ipxeInspect())
	mux.Handle("/ipxe/config", ipxeBoot(bootConfigs))
	return mux
}

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
func ipxeBoot(bootConfigs BootAdapter) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		params := req.URL.Query()
		attrs := MachineAttrs{UUID: params.Get("uuid")}
		bootConfig, err := bootConfigs.Get(attrs)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		var buf bytes.Buffer
		err = ipxeTemplate.Execute(&buf, bootConfig)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		buf.WriteTo(w)
	}
	return http.HandlerFunc(fn)
}
