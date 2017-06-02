package cli

import (
	"io/ioutil"
	"path/filepath"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// genericPutCmd creates and updates Generic templates.
var (
	genericPutCmd = &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create an Generic template",
		Long:  `Create an Generic template`,
		Run:   runGenericPutCmd,
	}
)

func init() {
	genericCmd.AddCommand(genericPutCmd)
	genericPutCmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "filename to use to create an Generic template")
	genericPutCmd.MarkFlagRequired("filename")
}

func runGenericPutCmd(cmd *cobra.Command, args []string) {
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
	req := &pb.GenericPutRequest{Name: filepath.Base(flagFilename), Config: config}
	_, err = client.Generic.GenericPut(context.TODO(), req)
	if err != nil {
		exitWithError(ExitError, err)
	}
}
