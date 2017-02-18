package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultServer = "localhost:60000"
)

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use:   "pop",
		Short: "A brief description of your application",
		Long: fmt.Sprintf(`Use pop to control a Pop server instance.
		
By default the server "%s" is used.
The client must authenticate with the server either via parameters specified through a POP_AUTH variable in the form
"username:password" or using the flags described below.`, DefaultServer),
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}
)

func init() {
	cobra.OnInitialize(viper.AutomaticEnv)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().String("auth", "", "username:password to use when authenticating with the Pop server")
	RootCmd.PersistentFlags().String("host", "", "specifies the server to connect to")
	RootCmd.PersistentFlags().Bool("json", false, "output JSON instead of YAML")

	viper.SetEnvPrefix("pop")

	viper.BindEnv("auth")
	viper.BindPFlag("auth", RootCmd.PersistentFlags().Lookup("auth"))

	viper.BindEnv("host")
	viper.BindPFlag("host", RootCmd.PersistentFlags().Lookup("host"))

	viper.BindPFlag("json", RootCmd.PersistentFlags().Lookup("json"))
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
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

func cl() *client.Client {
	return &client.Client{
		Credentials: credentials(),
	}
}

func credentials() creds.Credentials {
	user, pass := auth()
	host := viper.GetString("host")

	return creds.Credentials{
		Host:     host,
		Username: user,
		Password: pass,
	}
}

func getAllServerNames() []string {
	srvs, err := cl().Servers(context.Background())
	if err != nil {
		fail(err)
	}

	ret := make([]string, len(srvs))
	for i, srv := range srvs {
		ret[i] = srv.Name
	}

	return ret
}

func fail(v ...interface{}) {
	newV := append([]interface{}{"error: "}, v...)

	fmt.Fprintln(os.Stderr, newV...)
	os.Exit(1)
}

func failf(fstr string, params ...interface{}) {
	fail(fmt.Sprintf(fstr, params...))
}

func filter(v string) client.Filter {
	if len(v) > 3 && strings.HasPrefix(v, "id:") {
		return client.IDFilter(v[3:])
	}

	return client.NameFilter(v)
}

func results(out interface{}, err error) {
	step(out, err)
	client.FlushSessions()
	os.Exit(0)
}

func step(out interface{}, err error) {
	if err != nil {
		fail(err)
	}

	if out == nil {
		fmt.Println("ok")
		return
	}

	var bytes []byte
	if viper.GetBool("json") {
		bytes, err = json.MarshalIndent(out, "", "\t")
	} else {
		bytes, err = yaml.Marshal(out)
	}

	if err != nil {
		fail(err)
	}

	fmt.Println(string(bytes))
}
