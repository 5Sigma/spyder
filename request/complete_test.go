package request

import (
	"github.com/5sigma/spyder/config"
	"testing"
)

func TestOnComplete(t *testing.T) {
	ts := testServer(`
				{
					"inner": {
						"value": "1234567890"
					}
				}
			`)
	defer ts.Close()
	epConfig := endpointConfig(`
		{
			"url": "%s",
			"method": "post",
			"onComplete": ["variableSave"]
		}
	`, ts.URL)
	var variableSave = `
		$variables.set("var", JSON.parse($response.body).inner.value)
	`
	createFile("testdata/scripts/variableSave.js", variableSave)
	_, err := Do(epConfig)
	if err != nil {
		t.Fatalf("Error making request: %s", err)
	}
	val := config.GetVariable("var")
	if val != "1234567890" {
		t.Errorf("Config variable did not get set correctly in onComplete:\n%s", val)
	}
}
