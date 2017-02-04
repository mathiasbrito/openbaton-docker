package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Prints servers",
	Long:  `Prints a list of all the servers available on the server.
	
You can also specify a filter to query for a single server.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			results(cl().Servers(context.Background()))

		case 1:
			results(cl().Server(context.Background(), filter(args[0])))

		default:
			fail("too many parameters")
		}
	},
}

func init() {
	RootCmd.AddCommand(serversCmd)
}
