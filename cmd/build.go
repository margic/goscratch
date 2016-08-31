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
	"errors"
	"os"
	"os/exec"
	"path"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var outPath string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Execute go build on a package",
	Long: `Test executes test on a go package
		First gets the package then executes.
		Then executs go tests on the package.
		Writes results to mount`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("build called")
		err := runBuild()
		if err != nil {
			log.Error(err)
		}
		err = writeDockerFile()
		if err != nil {
			log.Error(err)
		}
	},
}

func runBuild() error {
	var pwd string
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = processRootFlags()
	if err != nil {
		return err
	}

	err = runGoGet()
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-v", "-o", path.Join(outPath, "goapp"))

	// get into source folder 	pwd, err = os.Getwd()
	if err != nil {
		return err
	}

	swd := path.Join(os.Getenv("GOPATH"), "src", goPkg)
	err = os.Chdir(swd)
	if err != nil {
		return err
	}
	defer os.Chdir(pwd)
	log.WithFields(
		log.Fields{
			"originalPath": pwd,
			"sourcepath":   swd,
		}).Debug("Path")
	// in source folder
	status, err := runCommand(cmd)
	if err == nil && status > 1 {
		return errors.New("Command returned non zero status")
	}
	return err
}

func writeDockerFile() error {
	log.Debug("Writing dockerfile")

	// create the report file
	f, err := os.Create(path.Join(outPath, "Dockerfile"))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	t, err := template.ParseFiles("template/Dockerfile")
	if err != nil {
		return err
	}

	t.Execute(w, data{
		Binary: "goapp",
	})

	w.Flush()
	return nil
}

type data struct {
	Binary string
}

func init() {
	RootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	buildCmd.Flags().StringVarP(&outPath, "output", "o", "/out", "Output folder for built assets default /out")

}
