package cmd

import (
	"fmt"
	"os"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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
	tw := newTabWriter(os.Stdout)
	defer tw.Flush()
	// legend
	fmt.Fprintf(tw, "ID\tNAME\tIGNITION\tCLOUD\tKERNEL\tINITRD\tCMDLINE\n")

	client := mustClientFromCmd(cmd)
	request := &pb.ProfileGetRequest{
		Id: args[0],
	}
	resp, err := client.Profiles.ProfileGet(context.TODO(), request)
	if err != nil {
		return
	}
	p := resp.Profile
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%#v\n", p.Id, p.Name, p.IgnitionId, p.CloudId, p.Boot.Kernel, p.Boot.Initrd, p.Boot.Cmdline)
}
