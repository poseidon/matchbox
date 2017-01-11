package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	pb "github.com/coreos/coreos-baremetal/matchbox/server/serverpb"
)

// groupListCmd lists Groups.
var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List machine groups",
	Long:  `List machine groups`,
	Run:   runGroupListCmd,
}

func init() {
	groupCmd.AddCommand(groupListCmd)
}

func runGroupListCmd(cmd *cobra.Command, args []string) {
	tw := newTabWriter(os.Stdout)
	defer tw.Flush()
	// legend
	fmt.Fprintf(tw, "ID\tGROUP NAME\tSELECTORS\tPROFILE\n")

	client := mustClientFromCmd(cmd)
	resp, err := client.Groups.GroupList(context.TODO(), &pb.GroupListRequest{})
	if err != nil {
		return
	}
	for _, group := range resp.Groups {
		fmt.Fprintf(tw, "%s\t%s\t%#v\t%s\n", group.Id, group.Name, group.Selector, group.Profile)
	}
}
