// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"github.com/mcilloni/openbaton-docker/pop/server"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "docker-popd",
	Short: "docker-popd server",
	Long: `docker-popd is a service that allows OpenBaton to orchestrate and deploy NFV on Docker containers.`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := loadConfig(); err != nil {
			log.WithError(err).Fatal("cannot load configuration file")
		}

		srv, err := server.New()
		if err != nil {
			log.WithError(err).Fatal("failure while launching popd")
		}

		if err := srv.Serve(); err != nil {
			log.WithError(err).Fatal("failure while running popd")
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		join := make(chan struct{})

		go func() {
			<-sigChan

			if err := srv.Close(); err != nil {
				log.WithError(err).Fatal("failure while stopping popd")
			}

			close(join)
		}()

		<-join
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "cfg", "", "config file (default is 'docker-popd.toml')")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" { // enable ability to specify config file via flag
		cfgFile = "docker-popd.toml"
	}
}

func loadConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	
	viper.AddConfigPath(wd)
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
