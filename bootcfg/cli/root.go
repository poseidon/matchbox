package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/coreos/coreos-baremetal/bootcfg/client"
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
		Endpoints []string
	}{}
)

func init() {
	RootCmd.PersistentFlags().StringSliceVar(&globalFlags.Endpoints, "endpoints", []string{"127.0.0.1:8081"}, "gRPC Endpoints")
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
	endpoints, err := cmd.Flags().GetStringSlice("endpoints")
	if err != nil {
		exitWithError(ExitError, err)
	}
	cfg := &client.Config{
		Endpoints: endpoints,
	}
	client, err := client.New(cfg)
	if err != nil {
		exitWithError(ExitBadConnection, err)
	}
	return client
}
