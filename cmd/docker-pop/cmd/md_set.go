package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets metadata values for a given ID",
	Long: `Sets the metadata for a server having the given ID, given in the format id {key=value}.
	
For instance, this command can be invoked as:
docker-pop md set <long UUID> key=val key1=val1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			failf("wrong number of arguments for md set: %d", len(args))
		}

		id := args[0]
		md := parseMd(args[1:])

		if err := cl().AddMetadata(context.Background(), id, md); err != nil {
			fail(err)
		}

		results(cl().FetchMetadata(context.Background(), id))
	},
}

func init() {
	mdCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func parseMd(args []string) map[string]string {
	ret := make(map[string]string, len(args))

	for _, arg := range args {
		splitted := strings.SplitN(arg, "=", 2)
		if len(splitted) != 2 {
			failf("invalid argument %s", arg)
		}

		ret[splitted[0]] = splitted[1]
	}

	return ret
}
