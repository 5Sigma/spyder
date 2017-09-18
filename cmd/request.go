package cmd

import (
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"github.com/5sigma/spyder/explorer"
	"github.com/5sigma/spyder/request"
	"github.com/5sigma/vox"
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
		if len(args) == 0 {
			cmd.Help()
			vox.Fatal("No endpoint specified")
		}
		epConfig, res := makeRequest(cmd, args[0])
		explorer.Start(args[0], epConfig, res)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("default", "d", false,
		"Use default values for all prompts")
}

func makeRequest(cmd *cobra.Command, endpointPath string) (*endpoint.EndpointConfig, *request.Response) {
	epConfig, err := endpoint.Load(path.Join(config.ProjectPath, "endpoints", endpointPath+".json"))
	if err != nil {
		vox.Fatal(err.Error())
	}
	for _, prompt := range epConfig.Prompts {
		useDefaults, _ := cmd.Flags().GetBool("default")
		if useDefaults {
			config.TempConfig.SetVariable(prompt.Name, prompt.DefaultValue)
		} else {
			value := vox.Prompt(prompt.Name, prompt.DefaultValue)
			config.TempConfig.SetVariable(prompt.Name, value)
		}
	}
	res, err := request.Do(epConfig)
	if err != nil {
		vox.Fatal(err.Error())
	}
	return epConfig, res
}
