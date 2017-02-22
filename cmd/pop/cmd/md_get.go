package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets metadata values for a given ID",
	Long: `Gets the metadata for a server having the given ID or name.
	
	pop md get nginx-cont
	or 
	pop md get id:uuid`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for md get: %d", len(args))
		}

		results(cl().FetchMetadata(context.Background(), filter(args[0])))
	},
}

func init() {
	mdCmd.AddCommand(getCmd)
}
