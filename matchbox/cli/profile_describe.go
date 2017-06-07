package cli

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// profileDescribeCmd describes a Profile.
var profileDescribeCmd = &cobra.Command{
	Use:   "describe PROFILE_ID",
	Short: "Describe a machine profile",
	Long:  `Describe a machine profile`,
	Run:   runProfileDescribeCmd,
}

func init() {
	profileCmd.AddCommand(profileDescribeCmd)
}

func runProfileDescribeCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}

	tw := newTabWriter(os.Stdout)
	defer tw.Flush()
	// legend
	fmt.Fprintf(tw, "ID\tNAME\tIGNITION\tCLOUD\tKERNEL\tINITRD\tARGS\n")

	client := mustClientFromCmd(cmd)
	request := &pb.ProfileGetRequest{
		Id: args[0],
	}
	resp, err := client.Profiles.ProfileGet(context.TODO(), request)
	if err != nil {
		return
	}
	p := resp.Profile
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", p.Id, p.Name, p.IgnitionId, p.CloudId, p.Boot.Kernel, p.Boot.Initrd, p.Boot.Args)
}
