package http

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
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

// logRequests logs HTTP requests.
func logRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("HTTP %s %v", req.Method, req.URL)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// selectGroup selects the Group whose selectors match the query parameters,
// adds the Group to the ctx, and calls the next handler. The next handler
// should handle a missing Group.
func selectGroup(srv server.Server, next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(req)
		// match machine request
		group, err := srv.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: attrs})
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
func selectProfile(srv server.Server, next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(req)
		// match machine request
		profile, err := srv.SelectProfile(ctx, &pb.SelectProfileRequest{Labels: attrs})
		if err == nil {
			// add the Profile to the ctx for the next handler
			ctx = withProfile(ctx, profile)
		}
		next.ServeHTTP(ctx, w, req)
	}
	return ContextHandlerFunc(fn)
}
