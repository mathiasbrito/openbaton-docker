package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// spawnCmd represents the spawn command
var spawnCmd = &cobra.Command{
	Use:   "spawn",
	Short: "Creates and then starts a server",
	Long: `spawn creates and starts a new server in a container, using the parameters provided as its arguments. 
Using spawn is equivalent to a create command followed by a start command.
	
Per example, spawn a new nginx container using the following parameters:

docker-pop spawn image=nginx name=nginx-cont
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pargs := parseArgs(args)

		results(
			cl().Spawn(
				context.Background(),
				pargs["name"],
				pargs["image"],
				pargs["flavour"],
				nil,
			),
		)
	},
}

func init() {
	RootCmd.AddCommand(spawnCmd)
}
