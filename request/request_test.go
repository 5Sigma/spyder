package request

import (
	"github.com/5Sigma/spyder/testhelper"
	"os"
	"testing"
)

var transform = `
	var body = JSON.parse($payload.get());
	body.inner.secondVar = 'two';
	$payload.set(JSON.stringify(body))
`

func TestMain(m *testing.M) {
	testhelper.Setup()
	testhelper.CreateFile("testdata/scripts/tansform.js", transform)
	retCode := m.Run()
	testhelper.Teardown()
	os.Exit(retCode)
}

func TestBadGetdata(t *testing.T) {
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
			"method": "GET",
			"data": {
				"outer": [
					{ "inner": "test" }
				]
			}
		}
	`, ts.URL)
	_, err := Do(epConfig)
	if err != nil {
		t.Fatalf("Error making request: %s", err)
	}
}
