package api

import (
	"bytes"
	"net/http"
	"text/template"

	"golang.org/x/net/context"
)

var grubTemplate = template.Must(template.New("GRUB2 config").Parse(`default=0
timeout=1
menuentry "CoreOS" {
echo "Loading kernel"
linuxefi "{{.Kernel}}"{{range $key, $value := .Cmdline}} {{if $value}}"{{$key}}={{$value}}"{{else}}"{{$key}}"{{end}}{{end}}
echo "Loading initrd"
initrdefi {{ range $element := .Initrd }}"{{$element}}" {{end}}
}
`))

// grubHandler returns a handler which renders a GRUB2 config for the
// requester.
func grubHandler() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		spec, err := specFromContext(ctx)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		var buf bytes.Buffer
		err = grubTemplate.Execute(&buf, spec.BootConfig)
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
