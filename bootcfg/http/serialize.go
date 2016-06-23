package http

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

// renderJSON encodes structs to JSON, writes the response to the
// ResponseWriter, and logs encoding errors.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		log.Errorf("error JSON encoding: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, js)
}

// writeJSON writes the given bytes with a JSON Content-Type.
func writeJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set(contentType, jsonContentType)
	_, err := w.Write(data)
	if err != nil {
		log.Errorf("error writing to response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func renderTemplate(w io.Writer, data interface{}, contents ...string) (err error) {
	tmpl := template.New("").Option("missingkey=error")
	for _, content := range contents {
		tmpl, err = tmpl.Parse(content)
		if err != nil {
			log.Errorf("error parsing template: %v", err)
			return err
		}
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Errorf("error rendering template: %v", err)
		return err
	}
	return nil
}
