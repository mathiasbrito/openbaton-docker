package cmd

import (
	"context"

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

docker-pop delete id:uuid`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for delete: %d", len(args))
		}

		results(nil, cl().Delete(context.Background(), filter(args[0])))
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
