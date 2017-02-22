package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a server",
	Long: `start starts a new server in a container, using the parameters provided as its arguments.
	
Per example, to start a server you have previously created, you may use the following invocation:

pop start nginx-cont
or
pop start id:uuid`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for start: %d", len(args))
		}

		results(cl().Start(context.Background(), filter(args[0])))
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
