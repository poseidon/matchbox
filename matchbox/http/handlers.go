package http

import (
	"fmt"
	"net/http"

	"github.com/poseidon/matchbox/matchbox/server"
	pb "github.com/poseidon/matchbox/matchbox/server/serverpb"
)

// homeHandler shows the server name for rooted requests. Otherwise, a 404 is
// returned.
func homeHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "matchbox\n")
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
func (s *Server) selectGroup(core server.Server, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		attrs := labelsFromRequest(s.logger, req)
		// match machine request
		group, err := core.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: attrs})
		if err == nil {
			// add the Group to the ctx for next handler
			ctx = withGroup(ctx, group)
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// selectProfile selects the Profile for the given query parameters, adds the
// Profile to the ctx, and calls the next handler. The next handler should
// handle a missing profile.
func (s *Server) selectProfile(core server.Server, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		attrs := labelsFromRequest(s.logger, req)
		// match machine request
		profile, err := core.SelectProfile(ctx, &pb.SelectProfileRequest{Labels: attrs})
		if err == nil {
			// add the Profile to the ctx for the next handler
			ctx = withProfile(ctx, profile)
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
