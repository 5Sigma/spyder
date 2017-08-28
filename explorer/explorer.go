package explorer

import (
	"fmt"
	"github.com/5sigma/gshell"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"github.com/5sigma/spyder/request"
	"github.com/5sigma/vox"
	"github.com/Jeffail/gabs"
	"github.com/dustin/go-humanize"
	"github.com/ttacon/chalk"
)

func Start(endpointPath string, epConfig *endpoint.EndpointConfig, res *request.Response) {
	shell := gshell.New()
	prompt := fmt.Sprintf("%s%s%s> ", chalk.Yellow, endpointPath, chalk.Reset)
	shell.Prompt = prompt
	if config.GetSetting("vimMode") == "true" {
		println("Enabling vim mode")
		shell.VimMode = true
	}

	shell.AddCommand(&gshell.Command{
		Name:        "response.body",
		Description: "Displays the content received from the server.",
		Call: func(sh *gshell.Shell, args []string) {
			if res.IsResponseJSON() {
				vox.PrintJSON(res.Content)
			} else {
				vox.Println(string(res.Content))
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "response",
		Description: "Displays various information about the response.",
		Call: func(sh *gshell.Shell, args []string) {
			vox.PrintProperty("Request Type", "HTTP")
			vox.PrintProperty("Request Time",
				fmt.Sprintf("%s", res.RequestTime))
			vox.PrintProperty("Response", string(res.Response.Status))
			vox.PrintProperty("Content Length", humanize.Bytes(uint64(res.Response.ContentLength)))
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "request",
		Description: "Displays various information about the request.",
		Call: func(sh *gshell.Shell, args []string) {
			vox.PrintProperty("Url", res.Request.URL.String())
			vox.PrintProperty("Content Length", humanize.Bytes(uint64(res.Request.ContentLength)))
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "request.body",
		Description: "Displays the post data sent in the request",
		Call: func(sh *gshell.Shell, args []string) {
			switch res.Request.Header.Get("Content-Type") {
			case "application/json":
				vox.PrintJSON(res.Payload)
			default:
				vox.Println(string(res.Payload))
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "get",
		Description: "Extracts a data out of a JSON response. A dot-notation JSON should be given as an argument.",
		Call: func(sh *gshell.Shell, args []string) {
			if len(args) == 0 {
				vox.Error("No path specified")
			} else {
				parsed, err := gabs.ParseJSON([]byte(res.Content))
				if err != nil {
					vox.Error(err.Error())
					return
				}
				value := parsed.Path(args[0]).String()
				if err != nil {
					vox.Error(err.Error())
					return
				}
				vox.PrintJSON([]byte(value))
			}
		},
	})

	shell.AddCommand(&gshell.Command{
		Name:        "refresh",
		Description: "Makes the request again and reloads all data.",
		Call: func(sh *gshell.Shell, args []string) {
			newRes, err := request.Do(epConfig)
			if err != nil {
				vox.Error(err.Error())
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
				vox.PrintProperty(key, value[0])
				if len(value) > 1 {
					for _, v := range value[1:] {
						vox.PrintProperty("", v)
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
					vox.PrintProperty(key, "")
				}
				if len(value) == 1 {
					vox.PrintProperty(key, value[0])
				}
				if len(value) > 1 {
					for _, v := range value[1:] {
						vox.PrintProperty("", v)
					}
				}
			}
		},
	})

	vox.Println("")
	shell.ProcessLine("response")
	vox.Println("")
	shell.ProcessLine("response.body")
	shell.Start()
}
