package cmd

import (
	"bytes"
	"github.com/5Sigma/spyder/config"
	"github.com/5Sigma/spyder/testhelper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testhelper.Setup()
	retCode := m.Run()
	testhelper.Teardown()
	os.Exit(retCode)
}

func TestPrompt(t *testing.T) {
	var (
		err error
	)
	ts := testhelper.RunTestServer(`
				{
					"inner": {
						"value": "1234567890"
					}
				}
			`)
	defer ts.Close()
	epConfig := testhelper.EndpointConfig(`
		{
			"url": "%s",
			"method": "POST",
			"prompts": {
				"val1": {
					"defaultValue": "aaa"
				}
			}
		}
	`, ts.URL)

	testhelper.CreateFile("testdata/endpoints/promptTest.json", epConfig.String())
	testhelper.AddPromptResponse("123")
	RootCmd.SetOutput(new(bytes.Buffer))
	RootCmd.SetArgs([]string{"req", "promptTest"})
	err = RootCmd.Execute()

	if err != nil {
		t.Fatalf("Error making request: %s", err)
	}

	val := config.GetVariable("val1")
	if val != "123" {
		t.Errorf("Config variable did not get set correctly from prompt:\n%s", val)
	}
}
