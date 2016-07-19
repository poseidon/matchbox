package http

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/mikeynap/coreos-baremetal/bootcfg/server"
	pb "github.com/mikeynap/coreos-baremetal/bootcfg/server/serverpb"
)

// requireGET requires requests to be an HTTP GET. Otherwise, it responds with
// a 405 status code.
func requireGET(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.Error(w, "only HTTP GET is supported", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// versionHandler shows the server name and version for root requests.
// Otherwise, a 404 is returned.
func versionHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "bootcfg\n")
	}
	return http.HandlerFunc(fn)
}

// logRequest logs HTTP requests.
func (s *Server) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		s.logger.Infof("HTTP %s %v", req.Method, req.URL)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// selectGroup selects the Group whose selectors match the query parameters,
// adds the Group to the ctx, and calls the next handler. The next handler
// should handle a missing Group.
func (s *Server) selectGroup(core server.Server, next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(s.logger, req)
		// match machine request
		group, err := core.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: attrs})
		if err == nil {
			// add the Group to the ctx for next handler
			ctx = withGroup(ctx, group)
		}
		next.ServeHTTP(ctx, w, req)
	}
	return ContextHandlerFunc(fn)
}

// selectProfile selects the Profile for the given query parameters, adds the
// Profile to the ctx, and calls the next handler. The next handler should
// handle a missing profile.
func (s *Server) selectProfile(core server.Server, next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(s.logger, req)
		// match machine request
		profile, err := core.SelectProfile(ctx, &pb.SelectProfileRequest{Labels: attrs})
		if err == nil {
			// add the Profile to the ctx for the next handler
			ctx = withProfile(ctx, profile)
		}
		next.ServeHTTP(ctx, w, req)
	}
	return ContextHandlerFunc(fn)
}
