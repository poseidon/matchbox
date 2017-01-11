package cli

import (
	"io/ioutil"
	"path/filepath"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/coreos-baremetal/matchbox/server/serverpb"
)

// ignitionPutCmd creates and updates Ignition templates.
var (
	ignitionPutCmd = &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create an Ignition template",
		Long:  `Create an Ignition template`,
		Run:   runIgnitionPutCmd,
	}
)

func init() {
	ignitionCmd.AddCommand(ignitionPutCmd)
	ignitionPutCmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "filename to use to create an Ignition template")
	ignitionPutCmd.MarkFlagRequired("filename")
}

func runIgnitionPutCmd(cmd *cobra.Command, args []string) {
	if len(flagFilename) == 0 {
		cmd.Help()
		return
	}
	if err := validateArgs(cmd, args); err != nil {
		return
	}

	client := mustClientFromCmd(cmd)
	config, err := ioutil.ReadFile(flagFilename)
	if err != nil {
		exitWithError(ExitError, err)
	}
	req := &pb.IgnitionPutRequest{Name: filepath.Base(flagFilename), Config: config}
	_, err = client.Ignition.IgnitionPut(context.TODO(), req)
	if err != nil {
		exitWithError(ExitError, err)
	}
}
