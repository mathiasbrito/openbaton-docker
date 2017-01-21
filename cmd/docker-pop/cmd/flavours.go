package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// flavoursCmd represents the flavours command
var flavoursCmd = &cobra.Command{
	Use:   "flavours",
	Short: "Prints flavours",
	Long: `Prints a list of all the flavours available on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		results(cl().Flavours(context.Background()))
	},
}

func init() {
	RootCmd.AddCommand(flavoursCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flavoursCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flavoursCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
