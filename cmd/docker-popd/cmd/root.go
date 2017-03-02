package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"github.com/mcilloni/openbaton-docker/docker-pop-server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	keepCont bool
	memprof  string
	verbose  bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "docker-popd",
	Short: "docker-popd server",
	Long:  `docker-popd is a service that allows OpenBaton to orchestrate and deploy NFV on Docker containers.`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := loadConfig(); err != nil {
			log.WithError(err).Fatal("cannot load configuration file")
		}

		srv, err := server.New()
		if err != nil {
			log.WithError(err).Fatal("failure while launching popd")
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		join := make(chan struct{})

		go func() {
			<-sigChan

			if err := srv.Close(); err != nil {
				srv.Logger.WithError(err).Fatal("failure while stopping popd")
			}

			close(join)
		}()

		if err := srv.Serve(); err != nil {
			srv.Logger.WithError(err).Fatal("failure while running popd")
		}

		<-join

		if memprof != "" {
			srv.Logger.WithField("memprof-file", memprof).Debug("attempting to dump memory profile file")

			f, err := os.Create(memprof)
			if err != nil {
				srv.Logger.WithError(err).Fatal("failure while creating profile file")
			}
			
			defer f.Close()
			
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				srv.Logger.WithError(err).Fatal("could not write memory profile")
			}

			srv.Logger.WithField("memprof-file", memprof).Debug("dumped memory profile file")
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "cfg", "", "config file (default is 'docker-popd.toml')")
	RootCmd.PersistentFlags().BoolVar(&keepCont, "keep-stopped", false, "keep Docker containers after they exit")
	RootCmd.PersistentFlags().StringVar(&memprof, "memprofile", "", "enable memory profiling, dumping the results to a given file (debug only)")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "output everything on the logs")

	viper.BindPFlag("keep-stopped", RootCmd.PersistentFlags().Lookup("keep-stopped"))
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
}

func loadConfig() error {
	if cfgFile == "" {
		return nil // will let the server self configure
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	viper.AddConfigPath(wd)
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
