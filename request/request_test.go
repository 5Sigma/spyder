package request

import (
	"fmt"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var transform = `
	var body = JSON.parse($payload.get());
	body.inner.secondVar = 'two';
	$payload.set(JSON.stringify(body)):
`

func endpointConfig(str, url string) *endpoint.EndpointConfig {
	createFile("testdata/endpoints/request.json",
		fmt.Sprintf(str, url))
	epConfig, _ := endpoint.Load(path.Join(config.ProjectPath,
		"endpoints", "request.json"))
	return epConfig
}

func TestMain(m *testing.M) {
	config.InMemory = true
	createFile("testdata/config.json", "")
	createFile("testdata/config.local.json", "")
	createFile("testdata/scripts/tansform.js", transform)
	config.ProjectPath = "testdata"

	retCode := m.Run()
	os.RemoveAll("testdata")
	os.Exit(retCode)
}

func folder(projectPath string, folder string) error {
	return os.MkdirAll(path.Join(projectPath, folder), os.ModePerm)
}

func createFile(fpath string, content string) error {
	os.MkdirAll(path.Dir(fpath), os.ModePerm)
	return ioutil.WriteFile(fpath, []byte(content), 0644)
}

func testServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
				{
					"inner": {
						"value": "1234567890"
					}
				}
			`)
		}))
}
