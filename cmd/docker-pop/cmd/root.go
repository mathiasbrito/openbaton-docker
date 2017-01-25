package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
	"github.com/ghodss/yaml"
)

const (
	DefaultServer = "localhost:60000"
)

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use:   "docker-pop",
		Short: "A brief description of your application",
		Long: fmt.Sprintf(`Use docker-pop to control a docker-popd instance.
		
By default the server "%s" is used.`, DefaultServer),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(viper.AutomaticEnv)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().String("auth", "", "username:password to log into popd with")
	RootCmd.PersistentFlags().String("host", "", "specifies the server to connect to")

	viper.SetEnvPrefix("pop")

	viper.BindEnv("auth")
	viper.BindPFlag("auth", RootCmd.PersistentFlags().Lookup("auth"))

	viper.BindEnv("host")
	viper.BindPFlag("host", RootCmd.PersistentFlags().Lookup("host"))
}

func fail(v ...interface{}) {
	newV := append([]interface{}{"error: "}, v...)

	fmt.Fprintln(os.Stderr, newV...)
	os.Exit(1)
}

func failf(fstr string, params ...interface{}) {
	fail(fmt.Sprintf(fstr, params...))
}

func cl() *client.Client {
	return &client.Client{
		Credentials: credentials(),
	}
}

func credentials() creds.Credentials {
	user, pass := auth()
	host := viper.GetString("host")

	return creds.Credentials{
		Host: host, 
		Username: user, 
		Password: pass,
	}
}

func auth() (string, string) {
	auth := viper.GetString("auth")
	if auth == "" {
		fail("no auth data provided")
	}

	splitted := strings.Split(auth, ":")
	if len(splitted) != 2 {
		failf("malformed auth string %s", auth)
	}

	return splitted[0], splitted[1]
}

func results(out interface{}, err error) {
	if err != nil {
		fail(err)
	}

	if out == nil {
		fmt.Println("ok")
		os.Exit(0)
	}

	bytes, err := yaml.Marshal(out)
	if err != nil {
		fail(err)
	}

	client.FlushSessions()

	fmt.Println(string(bytes))
	os.Exit(0)
}