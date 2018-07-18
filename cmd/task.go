package cmd

import (
	"github.com/5Sigma/spyder/config"
	"github.com/5Sigma/spyder/request"
	"github.com/5sigma/vox"
	"github.com/spf13/cobra"
	"path"
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Execute a task script.",
	Long: `Tasks are generalized scripts that can be executed. They are intended
to be used to perform multiple endpoint requests to easily make system changes
or automate actions.

Task scripts are Javascript files located in project/tasks folder.
When specifying a task you do not need to specify the extension.
`,
	Run: func(cmd *cobra.Command, args []string) {
		engine := request.NewTaskEngine()
		filepath := path.Join(config.ProjectPath, "tasks", args[0]+".js")
		err := engine.ExecuteFile(filepath)
		if err != nil {
			vox.Fatal(err.Error())
		}
	}}

func init() {
	RootCmd.AddCommand(taskCmd)
}
