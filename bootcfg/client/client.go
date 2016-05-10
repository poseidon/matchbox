package client

import (
	"errors"

	"google.golang.org/grpc"

	"github.com/coreos/coreos-baremetal/bootcfg/rpc/rpcpb"
)

var (
	errNoEndpoints = errors.New("client: No endpoints provided")
)

// Config configures a Client.go f
type Config struct {
	// List of endpoint URLs
	Endpoints []string
}

// Client provides a bootcfg client RPC session.
type Client struct {
	Groups   rpcpb.GroupsClient
	Profiles rpcpb.ProfilesClient
	Ignition rpcpb.IgnitionClient
	conn     *grpc.ClientConn
}

// New creates a new Client from the given Config.
func New(config *Config) (*Client, error) {
	if len(config.Endpoints) == 0 {
		return nil, errNoEndpoints
	}
	return newClient(config)
}

// Close closes the client's connections.
func (c *Client) Close() error {
	return c.conn.Close()
}

func newClient(config *Config) (*Client, error) {
	conn, err := retryDialer(config)
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:     conn,
		Groups:   rpcpb.NewGroupsClient(conn),
		Profiles: rpcpb.NewProfilesClient(conn),
		Ignition: rpcpb.NewIgnitionClient(conn),
	}
	return client, nil
}

// retryDialer attemps to Dial each endpoint in order to establish a
// connection.
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
