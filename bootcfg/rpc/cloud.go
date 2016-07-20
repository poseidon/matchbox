package rpc

import (
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/rpc/rpcpb"
	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// cloudServer takes a bootcfg Server and implements a gRPC IgnitionServer.
type cloudServer struct {
	srv server.Server
}

func newCloudServer(s server.Server) rpcpb.CloudServer {
	return &cloudServer{
		srv: s,
	}
}

func (s *cloudServer) CloudPut(ctx context.Context, req *pb.CloudPutRequest) (*pb.CloudPutResponse, error) {
	_, err := s.srv.CloudPut(ctx, req)
	return &pb.CloudPutResponse{}, grpcError(err)
}
