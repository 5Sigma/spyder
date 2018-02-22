package cmd

import (
	"fmt"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/vox"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
)

// endpointCmd represents the endpoint command
var endpointCmd = &cobra.Command{
	Use:   "endpoint [endpoint name]",
	Short: "Generate a new endpoint file",
	Long: `Create a new endpoint template at a given path.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			vox.Error("Specify an endpoint name such as auth/login")
			return
		}
		filePath := path.Join(config.ProjectPath, "endpoints", args[0]+".json")
		vox.PrintResult("Creating path", os.MkdirAll(filepath.Dir(filePath), os.ModePerm))
		vox.PrintResult("Creating file", writeFile(filePath, fmt.Sprintf(`
{
	"url": "%s",
	"method": "GET",
	"headers": {
		"Content-Type": "application/json"
	},
	"data": {
	},
	"definition": {
		"request": {
		},
		"response": {
		}
	}
}

@doc

Makes a call to the %s endpoint
`, args[0], args[0])))

	},
}

func init() {
	newCmd.AddCommand(endpointCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// endpointCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
