package cmd

import (
	"fmt"
	"os"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// profileListCmd lists Profiles.
var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List machine profiles",
	Long:  `List machine profiles`,
	Run:   runProfileListCmd,
}

func init() {
	profileCmd.AddCommand(profileListCmd)
}

func runProfileListCmd(cmd *cobra.Command, args []string) {
	tw := newTabWriter(os.Stdout)
	defer tw.Flush()
	// legend
	fmt.Fprintf(tw, "ID\tPROFILE NAME\tIGNITION\tCLOUD\n")

	client := mustClientFromCmd(cmd)
	resp, err := client.Profiles.ProfileList(context.TODO(), &pb.ProfileListRequest{})
	if err != nil {
		return
	}
	for _, profile := range resp.Profiles {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", profile.Id, profile.Name, profile.IgnitionId, profile.CloudId)
	}
}
