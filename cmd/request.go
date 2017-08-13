package cmd

import (
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"github.com/5sigma/spyder/explorer"
	"github.com/5sigma/spyder/output"
	"github.com/5sigma/spyder/request"
	"github.com/spf13/cobra"
	"path"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "request [endpoint]",
	Aliases: []string{"req", "get", "r"},
	Short:   "Make an endpoint request",
	Long: `Send a request to an endpoint using a given endpoint configration file.
You need only specifiy the relative path from the endpoints folder to the 
endpoint configuration. The file extension is also not needed.

For instance a configuration located at endpoints/sessions/auth.json can be
requested using:

$ spyder request sessions/auth
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := endpoint.Load(path.Join(config.ProjectPath, "endpoints", args[0]+".json"))
		if err != nil {
			output.PrintFatal(err)
		}

		res, err := request.Do(config)
		if err != nil {
			output.PrintFatal(err)
		}
		explorer.Start(args[0], config, res)

	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
