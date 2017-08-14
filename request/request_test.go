package request

import (
	"github.com/5sigma/spyder/testhelper"
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
