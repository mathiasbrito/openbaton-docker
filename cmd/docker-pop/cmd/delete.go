package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a server",
	Long: `delete deletes a stopped container, using the parameters provided as its arguments.
	
Per example, to delete a container named "nginx-cont" you have previously stopped or created, you may use the following invocation:

docker-pop delete nginx-cont

or

docker-pop delete id:uuid

The '*' symbol can be used to delete every server known to the daemon.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for delete: %d", len(args))
		}

		targets := args
		if args[0] == "*" {
			targets = getAllServerNames()
		}

		for _, target := range targets {
			fmt.Print(target, ": ")
			step(nil, cl().Delete(context.Background(), filter(target)))
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
