package cmd

import (
	"bytes"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/testhelper"
	"testing"
)

func TestTask(t *testing.T) {
	var (
		err error
	)

	ts := testhelper.RunEchoServer()
	defer ts.Close()
	epConfig := testhelper.EndpointConfig(`
		{
			"url": "%s",
			"method": "POST",
			"data": {
				"test": "test"
			}
		}
	`, ts.URL)

	testhelper.CreateFile("testdata/endpoints/ep.json", epConfig.String())
	script := `
		var res = $endpoint('ep', {"test": "123"});
		$variables.set('val1', res.body.test);
	`
	testhelper.CreateFile("testdata/tasks/task.js", script)
	RootCmd.SetOutput(new(bytes.Buffer))
	RootCmd.SetArgs([]string{"task", "task"})
	err = RootCmd.Execute()

	if err != nil {
		t.Fatalf("Command error: %s", err)
	}

	val := config.GetVariable("val1")
	if val != "123" {
		t.Errorf("Config variable did not get set correctly:\n%s", val)
	}
}
