package http

import (
	"net/http"

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
