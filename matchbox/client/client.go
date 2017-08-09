package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
)

var (
	errNoEndpoints = errors.New("client: No endpoints provided")
	errNoTLSConfig = errors.New("client: No TLS Config provided")
)

// Config configures a Client.
type Config struct {
	// List of endpoint URLs
	Endpoints []string
	// DialTimeout is the timeout for dialing a client connection
	DialTimeout time.Duration
	// Client TLS credentials
	TLS *tls.Config
}

// Client provides a matchbox client RPC session.
type Client struct {
	Groups   rpcpb.GroupsClient
	Profiles rpcpb.ProfilesClient
	Ignition rpcpb.IgnitionClient
	Generic  rpcpb.GenericClient
	conn     *grpc.ClientConn
}

// New creates a new Client from the given Config.
func New(config *Config) (*Client, error) {
	if len(config.Endpoints) == 0 {
		return nil, errNoEndpoints
	}
	for _, endpoint := range config.Endpoints {
		if _, _, err := net.SplitHostPort(endpoint); err != nil {
			return nil, fmt.Errorf("client: invalid host:port endpoint: %v", err)
		}
	}
	return newClient(config)
}

// Close closes the client's connections.
func (c *Client) Close() error {
	return c.conn.Close()
}

func newClient(config *Config) (*Client, error) {
	conn, err := dialEndpoints(config)
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:     conn,
		Groups:   rpcpb.NewGroupsClient(conn),
		Profiles: rpcpb.NewProfilesClient(conn),
		Ignition: rpcpb.NewIgnitionClient(conn),
		Generic:  rpcpb.NewGenericClient(conn),
	}
	return client, nil
}

// dialEndpoints attemps to Dial each endpoint in order to establish a
// connection.
func dialEndpoints(config *Config) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(config.DialTimeout),
	}
	if config.TLS != nil {
		creds := credentials.NewTLS(config.TLS)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		return nil, errNoTLSConfig
	}

	for _, endpoint := range config.Endpoints {
		conn, err = grpc.Dial(endpoint, opts...)
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
}
