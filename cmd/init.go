package cmd

import (
	"github.com/5sigma/vox"
	"github.com/spf13/cobra"
	"io/ioutil"
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
		vox.PrintResult("Created endpoints folder",
			createProjectFolder(projectPath, "endpoints"))
		vox.PrintResult("Created scripts folder",
			createProjectFolder(projectPath, "scripts"))
		vox.PrintResult("Created task folder",
			createProjectFolder(projectPath, "tasks"))
		vox.PrintResult("Create global config",
			writeFile("spyder.json", `
{
	"variables": {}
}
			`))
		vox.PrintResult("Create local config",
			writeFile("spyder.local.json", `
{
	"variables": {}
}
		`))
		vox.PrintResult("Created project", nil)
		vox.Println("\nProject files generated. If you version the project you should add 'spyder.local.json' to you're gitignore")
	},
}

func createProjectFolder(projectPath string, folder string) error {
	return os.MkdirAll(path.Join(projectPath, folder), os.ModePerm)
}

func writeFile(path string, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
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
