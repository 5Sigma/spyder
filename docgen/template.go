package docgen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/5Sigma/spyder/config"
	"github.com/5Sigma/spyder/endpoint"
	"github.com/5Sigma/spyder/sectionfile"
	"gopkg.in/russross/blackfriday.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type (
	EndpointData struct {
		Name            string
		Path            string
		Documentation   string
		BaseURL         string
		Method          string
		ResponseExample string
		RequestExample  string
		Headers         []endpoint.HeaderDefinition
		Params          string
	}

	APIData struct {
		ProjectName string
		BaseURL     string
		Endpoints   []*EndpointData
	}
)

func ProcessTemplate() (string, error) {
	data, err := getProjectData()
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("endpoint").Parse(IndexTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

func getProjectData() (*APIData, error) {
	var (
		pName   = config.GetSettingDefault("projectName", "API Documentation")
		apiData = &APIData{
			ProjectName: pName,
		}
		epFolder = path.Join(config.ProjectPath, "endpoints")
	)

	err := filepath.Walk(epFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".json" {
				newData, err := getEndpointData(path)
				if err != nil {
					return err
				}
				apiData.Endpoints = append(apiData.Endpoints, newData)
			}
			return nil
		},
	)
	return apiData, err
}

func getEndpointData(filename string) (*EndpointData, error) {
	ep, err := endpoint.Load(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]: %s", filename, err.Error()))
	}
	sf, err := sectionfile.Load(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]: %s", filename, err.Error()))
	}
	documentation := string(blackfriday.Run([]byte(sf.Contents("doc"))))
	documentation = strings.Replace(documentation, "\n", "", -1)
	documentation = strings.TrimSpace(documentation)
	epPath, _ := filepath.Rel(path.Join(config.ProjectPath, "endpoints"),
		strings.TrimSuffix(filename, filepath.Ext(filename)))
	name := epPath
	if ep.Name != "" {
		name = ep.Name
	}

	data := &EndpointData{
		Name:            name,
		Path:            epPath,
		BaseURL:         ep.Url,
		Documentation:   documentation,
		Method:          ep.Method,
		ResponseExample: ep.ExampleResponse(),
		RequestExample:  ep.ExampleRequest(),
		Headers:         ep.HeaderDefinitions(),
	}
	return data, nil
}
