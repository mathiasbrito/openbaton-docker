package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// flavoursCmd represents the flavours command
var flavoursCmd = &cobra.Command{
	Use:   "flavours",
	Short: "Prints flavours",
	Long: `Prints a list of all the flavours available on the server.
	
You can also specify a filter to query for a single flavour.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			results(cl().Flavours(context.Background()))

		case 1:
			results(cl().Flavour(context.Background(), filter(args[0])))

		default:
			fail("too many parameters")
		}
	},
}

func init() {
	RootCmd.AddCommand(flavoursCmd)
}
