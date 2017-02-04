package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Prints images",
	Long: `Prints a list of all the images available on the server.
	
You can also specify a filter to query for a single image.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			results(cl().Images(context.Background()))

		case 1:
			results(cl().Image(context.Background(), filter(args[0])))

		default:
			fail("too many parameters")
		}
	},
}

func init() {
	RootCmd.AddCommand(imagesCmd)
}
