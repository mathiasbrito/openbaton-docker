package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// networksCmd represents the networks command
var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Prints networks",
	Long: `Prints a list of all the networks available on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		results(cl().Networks(context.Background()))
	},
}

func init() {
	RootCmd.AddCommand(networksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// networksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// networksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
