package cmd

import (
	"context"
	"fmt"

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
docker-pop stop id:uuid

The '*' symbol can be used to stop every server known to the daemon.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for stop: %d", len(args))
		}

		targets := args
		if args[0] == "*" {
			targets = getAllServerNames()
		}

		for _, target := range targets {
			fmt.Print(target, ": ")
			step(nil, cl().Stop(context.Background(), filter(target)))
		}
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
