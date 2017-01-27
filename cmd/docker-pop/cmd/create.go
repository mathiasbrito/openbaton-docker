package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a server",
	Long: `create creates a new server in a container, using the parameters provided as its arguments.
	
Per example, create a new nginx container using the following parameters:

docker-pop create image=nginx name=nginx-cont
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pargs := parseArgs(args)

		results(
			cl().Create(
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
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func parseArgs(args []string) map[string]string {
	var name, image, flavour string

	for _, arg := range args {
		splitted := strings.SplitN(arg, "=", 2)
		if len(splitted) != 2 {
			failf("invalid argument %s", arg)
		}

		switch splitted[0] {
		case "image":
			image = splitted[1]

		case "flavour":
			flavour = splitted[1]

		case "name":
			name = splitted[1]

		default:
			failf("unknown argument %s", splitted[0])
		}
	}

	if image == "" {
		fail("no image chosen")
	}

	return map[string]string{
		"name":    name,
		"image":   image,
		"flavour": flavour,
	}
}
