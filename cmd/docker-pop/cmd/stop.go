package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a server",
	Long: `stop stops a running server, using the parameters provided as its arguments.
	
Per example, to stop a server you have previously started, you may use the following invocation:

docker-pop stop nginx-cont
or
docker-pop stop id:uuid`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for start: %d", len(args))
		}

		results(nil, cl().Stop(context.Background(), filter(args[0])))
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
