// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"bytes"
	"errors"
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// constants for use for configuration keys when using yaml file
// cfgLogging constant for the configuration key logging
const cfgLogging = "logging"

// cfgPackage constant for the configuration key package
const cfgPackage = "package"

var cfgFile string

// GoPkg the go package to act on
var goPkg string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "goscratch",
	Short: "Golang builder",
	Long: `Golang builder that build a docker image from a golang package.
Config can be provided by mounting a volume /etc/goscratch with a
goscratch.yaml config file.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	// StringVarP(p *string, name, shorthand string, value string, usage string)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /etc/goscratch/goscratch.yaml)")
	RootCmd.PersistentFlags().StringVarP(&goPkg, "package", "p", "", "golang package")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigType("yaml")
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
		log.Fatal("--config not working correctly")
	}

	// Set other path options
	viper.SetConfigName("goscratch")        // name of config file (without extension)
	viper.AddConfigPath("$HOME/.goscratch") // adding home directory as first search path
	viper.AddConfigPath("/etc/goscratch")   // adding etc as path. This path is mountable as a docker volume
	viper.AutomaticEnv()                    // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	// Set the logging level from config or environment variable override.
	ll := viper.GetString("logging")
	level, err := log.ParseLevel(ll)
	if err == nil {
		log.SetLevel(level)
	}
	log.Debug("Debug logging enabled.")
}

func processRootFlags() error {
	if goPkg == "" {
		// there is no package specified try the viperconfig
		goPkg = viper.GetString(cfgPackage)
		if goPkg == "" {
			return errors.New("no golang package specified")
		}
	}
	return nil
}

func runCommand(cmd *exec.Cmd) (status int, err error) {
	var waitStatus syscall.WaitStatus
	// set up output buffers
	// Stdout buffer
	// TODO fix this can't do this for tests with large output
	cmdOutput := &bytes.Buffer{}
	// Attach buffer to command
	cmd.Stdout = cmdOutput
	// Stderr buffer
	cmdError := &bytes.Buffer{}
	cmd.Stderr = cmdError

	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	}
	status = waitStatus.ExitStatus()
	log.WithField("ExitStatus", status).Debug("Command Exit Status")

	// output the commands stdout
	if cmdOutput.Len() > 0 {
		log.Infof("Output: %s\n", string(cmdOutput.Bytes()))
	}
	if cmdError.Len() > 0 {
		log.Infof("Error: %s\n", string(cmdError.Bytes()))
	}
	return status, err
}

func runGoGet() error {
	log.WithField("Package", goPkg).Debug("Running go get")
	cmd := exec.Command("go", "get", "-d", "-t", "-v", goPkg)

	status, err := runCommand(cmd)
	if err == nil && status > 1 {
		return errors.New("Command returned non zero status")
	}
	return err
}
