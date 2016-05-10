// Package client provides the bootcfg gRPC client.
//
// Create a bootcfg gRPC client using `client.New`:
//
//     cfg := &client.Config{
//       Endpoints: []string{"127.0.0.1:8081"},
//     }
//     client, err := client.New(cfg)
//     defer client.Close()
//
// Callers must Close the client after use.
//
package client
