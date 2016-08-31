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
	"errors"
	"os"
	"os/exec"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute go tests on a package",
	Long: `Test executes test on a go package
		First gets the package then executes.
		Then executs go tests on the package.
		Writes results to mount`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("test called")
		err := runTest()
		if err != nil {
			log.Error(err)
			cmd.Usage()
		}
	},
}

func runTest() error {
	var pwd string

	err := processRootFlags()
	if err != nil {
		return err
	}
	log.WithField("Package", goPkg).Debug("Running go get")
	cmd := exec.Command("go", "get", "-d", "-t", "-v", goPkg)

	status, err := runCommand(cmd)
	if err == nil && status > 1 {
		return errors.New("Command returned non zero status")
	}

	pwd, err = os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(path.Join(pwd, goPkg))
	log.WithField("pwd", pwd).Debug("Path")
	// in source folder
	//| go-junit-report > $CIRCLE_TEST_REPORTS/junit/test-results.xml
	cmd = exec.Command("go", "test", "-v", "./...", "|", "go-junit-report", ">", "/results/test-results.xml")

	status, err = runCommand(cmd)
	os.Chdir(pwd)
	if err == nil && status > 1 {
		return errors.New("Command returned non zero status")
	}

	log.Info("Test complete")
	return err
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
