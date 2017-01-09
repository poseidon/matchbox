package cli

import (
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage machine profiles",
	Long:  `List and describe machine profiles`,
}

func init() {
	RootCmd.AddCommand(profileCmd)
}
