package cmd

import (
	"fmt"
	"github.com/5Sigma/spyder/config"
	"github.com/5Sigma/spyder/docgen"
	"github.com/5sigma/vox"
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve API docs using a local webserver.",
	Long: `Spyder will use the endpoint definitions to generate API
documentation which will be served over http using a built in webserver.

The API documentation is generated using the request and response definitions,
and documentation markdown. For more information on how to provide information
for the documentation see the guide in the wiki:

https://github.com/5Sigma/spyder/wiki
`,
	Run: func(cmd *cobra.Command, args []string) {

		http.HandleFunc("/doc.js", func(w http.ResponseWriter, r *http.Request) {
			asset, _ := docgen.Asset("dist/doc.js")
			w.Write(asset)
		})
		http.HandleFunc("/doc.css", func(w http.ResponseWriter, r *http.Request) {
			asset, _ := docgen.Asset("dist/doc.css")
			w.Header().Set("Content-Type", "text/css")
			w.Write(asset)
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			println(r.URL.String())
			start := time.Now()
			tmpl, err := docgen.ProcessTemplate()
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, tmpl)
			}
			elapsed := time.Since(start)
			vox.Println("Documentation rebuild in ",
				vox.Yellow,
				fmt.Sprintf("%s", elapsed), vox.ResetColor)
		})

		port, _ := cmd.Flags().GetInt("port")
		hostStr := fmt.Sprintf("0.0.0.0:%d", port)
		vox.Println("Web server running on ", vox.Yellow, hostStr, vox.ResetColor)
		vox.Fatal(http.ListenAndServe(hostStr, nil))
	},
}

func init() {
	docsCmd.AddCommand(serveCmd)
	portStr := config.GetSettingDefault("docs.port", "3000")
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		portInt = 3000
	}
	serveCmd.PersistentFlags().Int("port", portInt, "Port to listen on")
}
