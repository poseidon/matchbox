package client

import (
	"google.golang.org/grpc"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// Config configures a Client.
type Config struct {
	// List of endpoint URLs
	Endpoints []string
}

// Client provides a bootcfg client RPC session.
type Client struct {
	Groups pb.GroupsClient
	conn *grpc.ClientConn
}

// New creates a new Client from the given Config.
func New(config *Config) (*Client, error) {
	return newClient(config)
}

func newClient(config *Config)  (*Client, error) {
	conn, err := retryDialer(config)
	if err != nil {
		return nil, err
	}
	client := &Client{
		Groups: pb.NewGroupsClient(conn),
		conn: conn,
	}
	return client, nil
}

// retryDialer attemps to Dial each endpoint until a client connection
// is established.
func retryDialer(config *Config) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	var err error
	for _, endp := range config.Endpoints {
		conn, err := grpc.Dial(endp, opts...)
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
}
