package explorer

import (
	"errors"
	"fmt"
	"github.com/5sigma/gshell"
	"github.com/5sigma/spyder/endpoint"
	"github.com/5sigma/spyder/output"
	"github.com/5sigma/spyder/request"
	"github.com/Jeffail/gabs"
	"github.com/dustin/go-humanize"
	"github.com/ttacon/chalk"
)

func Start(endpointPath string, config *endpoint.EndpointConfig, res *request.Response) {
	shell := gshell.New()
	prompt := fmt.Sprintf("%s%s%s> ", chalk.Yellow, endpointPath, chalk.Reset)
	shell.Prompt = prompt

	shell.AddCommand(&gshell.Command{
		Name:        "response.body",
		Description: "Displays the content received from the server.",
		Call: func(sh *gshell.Shell, args []string) {
			if res.IsResponseJSON() {
				output.PrintJson(res.Content)
			} else {
				fmt.Println(string(res.Content))
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "response",
		Description: "Displays various information about the response.",
		Call: func(sh *gshell.Shell, args []string) {
			output.PrintProperty("Request Type", "HTTP")
			output.PrintProperty("Request Time",
				fmt.Sprintf("%s", res.RequestTime))
			output.PrintProperty("Response", string(res.Response.Status))
			output.PrintProperty("Content Length", humanize.Bytes(uint64(res.Response.ContentLength)))
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "request",
		Description: "Displays various information about the request.",
		Call: func(sh *gshell.Shell, args []string) {
			output.PrintProperty("Url", res.Request.URL.String())
			output.PrintProperty("Content Length", humanize.Bytes(uint64(res.Request.ContentLength)))
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "request.body",
		Description: "Displays the post data sent in the request",
		Call: func(sh *gshell.Shell, args []string) {
			switch res.Request.Header.Get("Content-Type") {
			case "application/json":
				output.PrintJson(res.Payload)
			default:
				fmt.Println(string(res.Payload))
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "get",
		Description: "Extracts a data out of a JSON response. A dot-notation JSON should be given as an argument.",
		Call: func(sh *gshell.Shell, args []string) {
			if len(args) == 0 {
				output.PrintError(errors.New("No path specified"))
			} else {
				parsed, err := gabs.ParseJSON([]byte(res.Content))
				if err != nil {
					output.PrintError(err)
					return
				}
				value := parsed.Path(args[0]).String()
				if err != nil {
					output.PrintError(err)
					return
				}
				output.PrintJson([]byte(value))
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "refresh",
		Description: "Makes the request again and reloads all data.",
		Call: func(sh *gshell.Shell, args []string) {
			newRes, err := request.Do(config)
			if err != nil {
				output.PrintError(err)
			}
			res = newRes
			sh.ProcessLine("response")
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "response.headers",
		Description: "Displays the headers recieved from the server",
		Call: func(sh *gshell.Shell, args []string) {
			for key, value := range res.Response.Header {
				output.PrintProperty(key, value[0])
				if len(value) > 1 {
					for _, v := range value[1:] {
						output.PrintProperty("", v)
					}
				}
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "request.headers",
		Description: "Displays the headers sent to the server",
		Call: func(sh *gshell.Shell, args []string) {
			for key, value := range res.Request.Header {
				if len(value) == 0 {
					output.PrintProperty(key, "")
				}
				if len(value) == 1 {
					output.PrintProperty(key, value[0])
				}
				if len(value) > 1 {
					for _, v := range value[1:] {
						output.PrintProperty("", v)
					}
				}
			}
		},
	})

	fmt.Println("")
	shell.ProcessLine("response")
	fmt.Println("")
	shell.ProcessLine("response.body")
	shell.Start()
}
