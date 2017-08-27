package request

import (
	"github.com/5sigma/spyder/testhelper"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransform(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			bytes, _ := ioutil.ReadAll(r.Body)
			json, _ := gabs.ParseJSON(bytes)
			val, _ := json.Path("var").Data().(string)
			if val != "hello" {
				t.Errorf("Request was not transformed: %s", val)
			}
		}))
	defer ts.Close()
	epConfig := testhelper.EndpointConfig(`
		{
			"url": "%s",
			"method": "post",
			"transform": ["transform"],
			"data": {
				"var": "123"
			}
		}
	`, ts.URL)
	transform := `
		$request.body.var = "hello";
		$request.setBody($request.body);
	`
	testhelper.CreateFile("testdata/scripts/transform.js", transform)
	_, err := Do(epConfig)
	if err != nil {
		t.Errorf("Request error: %s", err.Error())
	}
}
