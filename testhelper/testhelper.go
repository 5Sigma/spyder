package testhelper

import (
	"fmt"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"github.com/5sigma/spyder/output"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
)

var PromptStream *os.File

func Setup() {
	config.InMemory = true
	output.SupressOutput = true
	CreateFile("testdata/config.json", "")
	CreateFile("testdata/config.local.json", "")
	config.ProjectPath = "testdata"
	PromptStream, _ = ioutil.TempFile("", "")
	output.PromptStream = PromptStream
}

func Teardown() {
	os.RemoveAll("testdata")
}

func Folder(projectPath string, folder string) error {
	return os.MkdirAll(path.Join(projectPath, folder), os.ModePerm)
}

func CreateFile(fpath string, content string) error {
	os.MkdirAll(path.Dir(fpath), os.ModePerm)
	return ioutil.WriteFile(fpath, []byte(content), 0644)
}

func RunTestServer(response string) *httptest.Server {
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

func EndpointConfig(str, url string) *endpoint.EndpointConfig {
	CreateFile("testdata/endpoints/request.json",
		fmt.Sprintf(str, url))
	epConfig, _ := endpoint.Load(path.Join(config.ProjectPath,
		"endpoints", "request.json"))
	return epConfig
}

func AddPromptResponse(str string) {
	io.WriteString(PromptStream, fmt.Sprintf("%s\n", str))
	PromptStream.Seek(0, os.SEEK_SET)
}
