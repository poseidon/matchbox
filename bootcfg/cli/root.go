package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mikeynap/coreos-baremetal/bootcfg/client"
	"github.com/mikeynap/coreos-baremetal/bootcfg/tlsutil"
)

var (
	// RootCmd is the base bootcmd command.
	RootCmd = &cobra.Command{
		Use:   "bootcmd",
		Short: "A command line client for the bootcfg service.",
		Long: `A CLI for the bootcfg Service

To get help about a resource or command, run "bootcmd help resource"`,
	}

	// globalFlags can be set for any subcommand.
	globalFlags = struct {
		endpoints []string
		caFile    string
		certFile  string
		keyFile   string
	}{}
)

func init() {
	RootCmd.PersistentFlags().StringSliceVar(&globalFlags.endpoints, "endpoints", []string{"127.0.0.1:8081"}, "gRPC Endpoints")
	// gRPC TLS Server Verification
	RootCmd.PersistentFlags().StringVar(&globalFlags.caFile, "ca-file", "/etc/bootcfg/ca.crt", "Path to the CA bundle to verify certificates of TLS servers")
	// gRPC TLS Client Authentication
	RootCmd.PersistentFlags().StringVar(&globalFlags.certFile, "cert-file", "/etc/bootcfg/client.crt", "Path to the client TLS certificate file")
	RootCmd.PersistentFlags().StringVar(&globalFlags.keyFile, "key-file", "/etc/bootcfg/client.key", "Path to the client TLS key file")
	cobra.EnablePrefixMatching = true
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// mustClientFromCmd returns a gRPC client or exits.
func mustClientFromCmd(cmd *cobra.Command) *client.Client {
	endpoints := endpointsFromCmd(cmd)
	tlsinfo := tlsInfoFromCmd(cmd)

	// client config
	tlscfg, err := tlsinfo.ClientConfig()
	if err != nil {
		exitWithError(ExitBadArgs, err)
	}
	cfg := &client.Config{
		Endpoints: endpoints,
		TLS:       tlscfg,
	}

	// gRPC client
	client, err := client.New(cfg)
	if err != nil {
		exitWithError(ExitBadConnection, err)
	}
	return client
}

// endpointsFromCmd returns the endpoint arguments.
func endpointsFromCmd(cmd *cobra.Command) []string {
	endpoints, err := cmd.Flags().GetStringSlice("endpoints")
	if err != nil {
		exitWithError(ExitBadArgs, err)
	}
	return endpoints
}

// tlsInfoFromCmd collects TLS arguments and returns a TLSInfo struct.
func tlsInfoFromCmd(cmd *cobra.Command) *tlsutil.TLSInfo {
	caFile, err := cmd.Flags().GetString("ca-file")
	if err != nil {
		exitWithError(ExitBadArgs, err)
	}
	certFile, err := cmd.Flags().GetString("cert-file")
	if err != nil {
		exitWithError(ExitBadArgs, err)
	}
	keyFile, err := cmd.Flags().GetString("key-file")
	if err != nil {
		exitWithError(ExitBadArgs, err)
	}

	return &tlsutil.TLSInfo{
		CAFile:   caFile,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}
