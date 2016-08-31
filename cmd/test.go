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
	"bufio"
	"io"
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
	var results = "/results/test-results.xml"

	err := processRootFlags()
	if err != nil {
		return err
	}

	err = runGoGet()
	if err != nil {
		return err
	}
	pwd, err = os.Getwd()
	if err != nil {
		return err
	}
	swd := path.Join(pwd, goPkg)
	os.Chdir(swd)
	defer os.Chdir(pwd)
	log.WithFields(
		log.Fields{
			"originalPath": pwd,
			"sourcepath":   swd,
		}).Debug("Path")
	// in source folder
	//| go-junit-report > $CIRCLE_TEST_REPORTS/junit/test-results.xml
	test := exec.Command("go", "test", "-v", "./...")
	report := exec.Command("go-junit-report")

	r, w := io.Pipe()

	// pipe test out to report in
	test.Stdout = w
	report.Stdin = r

	// create the report file
	f, err := os.Create(results)
	if err != nil {
		return err
	}
	defer f.Close()

	reportwriter := bufio.NewWriter(f)
	// attach file out to report command output
	report.Stdout = reportwriter

	// run commnads and wait for completion
	test.Start()
	report.Start()
	test.Wait()
	w.Close()
	report.Wait()

	reportwriter.Flush()
	log.WithField("Output", results).Info("Test complete")
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
