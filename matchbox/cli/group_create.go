package cli

import (
	"io/ioutil"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// groupPutCmd creates and updates Groups.
var (
	groupPutCmd = &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a machine group",
		Long:  `Create a machine group`,
		Run:   runGroupPutCmd,
	}
)

func init() {
	groupCmd.AddCommand(groupPutCmd)
	groupPutCmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "filename to use to create a Group")
	groupPutCmd.MarkFlagRequired("filename")
	groupPutCmd.MarkFlagFilename("filename", "json")
}

func runGroupPutCmd(cmd *cobra.Command, args []string) {
	if len(flagFilename) == 0 {
		cmd.Help()
		return
	}
	if err := validateArgs(cmd, args); err != nil {
		return
	}

	client := mustClientFromCmd(cmd)
	group, err := loadGroup(flagFilename)
	if err != nil {
		exitWithError(ExitError, err)
	}
	req := &pb.GroupPutRequest{Group: group}
	_, err = client.Groups.GroupPut(context.TODO(), req)
	if err != nil {
		exitWithError(ExitError, err)
	}
}

func loadGroup(filename string) (*storagepb.Group, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseGroup(data)
}
