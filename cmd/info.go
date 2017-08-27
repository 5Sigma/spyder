package cmd

import (
	"fmt"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/output"
	"github.com/Jeffail/gabs"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display project information",
	Long:  `Display information about Spyder and the current project.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintProperty("Project Path", config.ProjectPath)
		epFolder := path.Join(config.ProjectPath, "endpoints")
		var epCount int = 0
		epErrors := make(map[string]string)
		filepath.Walk(epFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				output.PrintFatal(err)
				return err
			}
			if filepath.Ext(path) == ".json" {
				epCount++
				_, err := gabs.ParseJSONFile(path)
				if err != nil {
					epErrors[path] = err.Error()
				}
			}
			return nil
		})
		output.PrintProperty("Endpoints", fmt.Sprintf("%d", epCount))
		if len(epErrors) > 0 {
			output.PrintErrorStr("Endpoint Errors:")
			for k, v := range epErrors {
				output.Println(k)
				output.PrintErrorStr(v)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
