package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// networksCmd represents the networks command
var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Prints networks",
	Long: `Prints a list of all the networks available on the server.
	
You can also specify a filter to query for a single network.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			results(cl().Networks(context.Background()))

		case 1:
			results(cl().Network(context.Background(), filter(args[0])))

		default:
			fail("too many parameters")
		}
	},
}

func init() {
	RootCmd.AddCommand(networksCmd)
}
