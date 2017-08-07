package cmd

import (
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
	Aliases: []string{"req", "get"},
	Short:   "Make an endpoint request",
	Long: `Send a request to an endpoint using a given endpoint configration file
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := endpoint.Load(path.Join("endpoints", args[0]+".json"))
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
	runCmd.Flags().BoolP("interactive", "i", false, "Run in interactive mode")
}
