package cmd

import (
	"fmt"
	"github.com/5sigma/vox"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:     "validate [endpoint]",
	Aliases: []string{"v"},
	Short:   "Validates an endpoint's response.",
	Long: `Validate the response from an endoint against a defined schema
definition. In order to validate the endpoint the endpoint's configuration 
needs to have an exepcted schema definition. This is located at
definition.response in the configuration file:

{
	...
	"definition": {
		"response": {
			"type": "object",
			"properties": {
				"id": {
					"type": "number"
				}
			}
		}
	}
}

The validation syntax is based on JsonSchema Draft 4:
https://tools.ietf.org/html/draft-zyp-json-schema-04
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			vox.Fatal("No endpoint specified")
		}
		epConfig, res := makeRequest(cmd, args[0])
		result, err := epConfig.ValidateResponse(string(res.Content))
		if err != nil {
			vox.Fatal(err.Error())
		}
		if result.Valid() {
			vox.Printlnc(vox.Green, fmt.Sprintf("[%s] valid response", args[0]))
		} else {
			vox.Error("The response was invalid:\n")
			for _, desc := range result.Errors() {
				vox.Errorf("- %s\n", desc)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
