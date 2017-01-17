package cli

import (
	"io/ioutil"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// profilePutCmd creates and updates Profiles.
var (
	profilePutCmd = &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a machine profile",
		Long:  `Create a machine profile`,
		Run:   runProfilePutCmd,
	}
	flagFilename string
)

func init() {
	profileCmd.AddCommand(profilePutCmd)
	profilePutCmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "filename to use to create a Profile")
	profilePutCmd.MarkFlagRequired("filename")
	profilePutCmd.MarkFlagFilename("filename", "json")
}

func runProfilePutCmd(cmd *cobra.Command, args []string) {
	if len(flagFilename) == 0 {
		cmd.Help()
		return
	}
	if err := validateArgs(cmd, args); err != nil {
		return
	}

	client := mustClientFromCmd(cmd)
	profile, err := loadProfile(flagFilename)
	if err != nil {
		exitWithError(ExitError, err)
	}
	req := &pb.ProfilePutRequest{Profile: profile}
	_, err = client.Profiles.ProfilePut(context.TODO(), req)
	if err != nil {
		exitWithError(ExitError, err)
	}
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return usageError(cmd, "Unexpected args: %v", args)
	}
	return nil
}

func loadProfile(filename string) (*storagepb.Profile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseProfile(data)
}
