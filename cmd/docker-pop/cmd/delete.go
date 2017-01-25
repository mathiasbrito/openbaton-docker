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

docker-pop delete nginx-cont`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			failf("wrong number of arguments for delete: %d", len(args))
		}

		results(nil, cl().Delete(context.Background(), args[0]))
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}