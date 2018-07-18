package request

import (
	"github.com/5Sigma/spyder/config"
	"github.com/5Sigma/spyder/testhelper"
	"testing"
)

func TestOnComplete(t *testing.T) {
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
			"method": "post",
			"onComplete": ["variableSave"]
		}
	`, ts.URL)
	var variableSave = `
		$variables.set("var", $response.body.inner.value)
	`
	testhelper.CreateFile("testdata/scripts/variableSave.js", variableSave)
	_, err := Do(epConfig)
	if err != nil {
		t.Fatalf("Error making request: %s", err)
	}
	val := config.GetVariable("var")
	if val != "1234567890" {
		t.Errorf("Config variable did not get set correctly in onComplete:\n%s", val)
	}
}
