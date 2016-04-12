package cli

import (
	"fmt"
	"os"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// groupDescribeCmd describes a Group.
var groupDescribeCmd = &cobra.Command{
	Use:   "describe GROUP_ID",
	Short: "Describe a machine group",
	Long:  `Describe a machine group`,
	Run:   runGroupDescribeCmd,
}

func init() {
	groupCmd.AddCommand(groupDescribeCmd)
}

func runGroupDescribeCmd(cmd *cobra.Command, args []string) {
	tw := newTabWriter(os.Stdout)
	defer tw.Flush()
	// legend
	fmt.Fprintf(tw, "ID\tNAME\tSELECTORS\tPROFILE\tMETADATA\n")

	client := mustClientFromCmd(cmd)
	request := &pb.GroupGetRequest{
		Id: args[0],
	}
	resp, err := client.Groups.GroupGet(context.TODO(), request)
	if err != nil {
		return
	}
	g := resp.Group
	fmt.Fprintf(tw, "%s\t%s\t%s\t%#v\t%s\n", g.Id, g.Name, g.Requirements, g.Profile, g.Metadata)
}
