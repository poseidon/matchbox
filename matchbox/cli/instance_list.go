package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// instanceListCmd represents the instance list command
var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List observed machine instances",
	Long:  `List observed machine instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init() {
	instanceCmd.AddCommand(instanceListCmd)
}
