package http

import (
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

// ContextHandler defines a handler which receives a passed context.Context
// with the standard ResponseWriter and Request.
type ContextHandler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

// ContextHandlerFunc type is an adapter to allow the use of an ordinary
// function as a ContextHandler. If f is a function with the correct
// signature, ContextHandlerFunc(f) is a ContextHandler that calls f.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTP calls the function f(ctx, w, req).
func (f ContextHandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	f(ctx, w, req)
}

// handler wraps a ContextHandler to implement the http.Handler interface for
// compatability with ServeMux and middlewares.
//
// Middleswares which do not pass a ctx break the chain so place them before
// or after chains of ContextHandlers.
type handler struct {
	ctx     context.Context
	handler ContextHandler
}

// NewHandler returns an http.Handler which wraps the given ContextHandler
// and creates a background context.Context.
func NewHandler(h ContextHandler) http.Handler {
	return &handler{
		ctx:     context.Background(),
		handler: h,
	}
}

// ServeHTTP lets handler implement the http.Handler interface.
func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(h.ctx, w, req)
}

// labelsFromRequest returns request query parameters.
func labelsFromRequest(req *http.Request) map[string]string {
	values := req.URL.Query()
	labels := map[string]string{}
	for key := range values {
		switch strings.ToLower(key) {
		case "mac":
			// set mac if and only if it parses
			if hw, err := parseMAC(values.Get(key)); err == nil {
				labels[key] = hw.String()
			}
		default:
			// matchers don't use multi-value keys, drop later values
			labels[key] = values.Get(key)
		}
	}
	return labels
}

// parseMAC wraps net.ParseMAC with logging.
func parseMAC(s string) (net.HardwareAddr, error) {
	macAddr, err := net.ParseMAC(s)
	if err != nil {
		// invalid MAC arguments may be common
		log.Debugf("error parsing MAC address: %s", err)
		return nil, err
	}
	return macAddr, err
}
