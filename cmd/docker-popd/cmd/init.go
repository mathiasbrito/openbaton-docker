package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/mcilloni/openbaton-docker/pop/server"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises a new configuration file",
	Long: `init initialises a new configuration file called docker-popd.toml in the current directory.
	You will be asked for an username and a password.
	
	Any existing config file will be overwritten.`,

	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read password: %v\n", err)
			os.Exit(1)
		}

		password := string(bytePassword)

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		user, err := server.NewUser(username, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot hash password: %v\n", err)
			os.Exit(1)
		}

		cfg := server.DefaultConfig
		cfg.Users = server.Users{username: user}

		if err := cfg.StoreFile(cfgFile, true); err != nil {
			fmt.Fprintf(os.Stderr, "cannot create config file \"%s\": %v\n", cfgFile, err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
