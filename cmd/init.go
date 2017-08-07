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
	"github.com/5sigma/spyder/output"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new project",
	Long: `Sets up a new project by creating standard folder structures and files. By
default it generates the project in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		var projectPath = "."
		if len(args) > 0 {
			projectPath = args[0]
		}
		output.PrintResult("Created endpoints folder",
			createProjectFolder(projectPath, "endpoints"))
		output.PrintResult("Created config folder",
			createProjectFolder(projectPath, "config"))
		output.PrintResult("Created scripts folder",
			createProjectFolder(projectPath, "scripts"))
		output.PrintResult("Created project", nil)
	},
}

func createProjectFolder(projectPath string, folder string) error {
	return os.MkdirAll(path.Join(projectPath, folder), os.ModePerm)
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
