package cmd

import (
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/docgen"
	"github.com/5sigma/vox"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"gen", "generate", "b", "g"},
	Short:   "Generate a static site for API documentation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		vox.PrintResult("Setting up build folder", createProjectFolder(config.ProjectPath, "docs"))
		template, err := docgen.ProcessTemplate()
		if err != nil {
			vox.PrintResult("Building html template", err)
		} else {
			err := ioutil.WriteFile(path.Join(config.ProjectPath, "docs", "index.html"), []byte(template), 660)
			vox.PrintResult("Building html template", err)
		}
		jsAsset, _ := docgen.Asset("dist/doc.js")
		vox.PrintResult("Building Javascript assets ",
			ioutil.WriteFile(path.Join(config.ProjectPath, "docs", "doc.js"), jsAsset, 660))
		cssAsset, _ := docgen.Asset("dist/doc.css")
		vox.PrintResult("Building CSS assets",
			ioutil.WriteFile(path.Join(config.ProjectPath, "docs", "doc.css"), cssAsset, 660))
	},
}

func init() {
	docsCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
