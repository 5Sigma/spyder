package cmd

import (
	"errors"
	"fmt"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"github.com/5sigma/spyder/output"
	"github.com/5sigma/spyder/request"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"path"
	"time"
)

// hammerCmd represents the hammer command
var hammerCmd = &cobra.Command{
	Use:   "hammer",
	Short: "Makes an endpoint request a number of times rapidly.",
	Long: `Make a number of request to an endpoint very rapidly and record the request timing.  The hammer command expects an endpoint to be passed in the same manner as the 'request' command:

spyder hammer --count 100 myEndpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			count      int
			totalTime  time.Duration
			maxTime    time.Duration
			minTime    time.Duration
			totalBytes int64
		)

		count, _ = cmd.Flags().GetInt("count")
		if len(args) == 0 {
			cmd.Help()
			output.PrintFatal(errors.New("No endpoint specified"))
		}

		configPath := path.Join(config.ProjectPath, "endpoints", args[0]+".json")
		epConfig, err := endpoint.Load(configPath)
		if err != nil {
			output.PrintFatal(err)
		}

		for _, prompt := range epConfig.Prompts {
			useDefaults, _ := cmd.Flags().GetBool("default")
			if useDefaults {
				config.TempConfig.SetVariable(prompt.Name, prompt.DefaultValue)
			} else {
				value := output.Prompt(prompt.Name, prompt.DefaultValue)
				config.TempConfig.SetVariable(prompt.Name, value)
			}
		}

		bar := output.NewProgress(count)
		for i := 0; i <= count; i++ {
			res, err := request.Do(epConfig)
			if err != nil {
				output.PrintFatal(err)
			}
			totalTime += res.RequestTime
			bar.Inc()
			if minTime == 0 {
				minTime = res.RequestTime
			}
			if res.RequestTime > maxTime {
				maxTime = res.RequestTime
			}
			if res.RequestTime < minTime {
				minTime = res.RequestTime
			}
			totalBytes += res.Response.ContentLength
		}

		avgTime := totalTime / time.Duration(count)

		output.PrintProperty("Number of requests", fmt.Sprintf("%d", count))
		output.PrintProperty("Average time", fmt.Sprintf("%s", avgTime))
		output.PrintProperty("Fastest", fmt.Sprintf("%s", minTime))
		output.PrintProperty("Slowest", fmt.Sprintf("%s", maxTime))
		output.PrintProperty("Total data transfer",
			humanize.Bytes(uint64(totalBytes)))
	},
}

func init() {
	RootCmd.AddCommand(hammerCmd)
	hammerCmd.PersistentFlags().Int("count", 100, "Request count")
	hammerCmd.Flags().BoolP("default", "d", false,
		"Use default values for all prompts")
}
