package cmd

import (

	"github.com/spf13/cobra"
)

// mdCmd represents the md command
var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "Manages metadata",
	Long: `Manages metadata for servers`,
}

func init() {
	RootCmd.AddCommand(mdCmd)
}
