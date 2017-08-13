package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ttacon/chalk"
	"os"
	"regexp"
	"strings"
)

func PrintResult(desc string, err error) {
	resultColor := chalk.Red
	resultText := "FAIL"
	if err == nil {
		resultColor = chalk.Green
		resultText = "OK"
	}
	desc += strings.Repeat(" ", 50-len(desc))
	fmt.Println(chalk.White, desc, chalk.Yellow, "[", resultColor, resultText,
		chalk.Yellow, "]", chalk.Reset)
}

func PrintError(err error) {
	fmt.Println(chalk.Red, err.Error(), chalk.Reset)
}

func PrintFatal(err error) {
	PrintError(err)
	os.Exit(1)
}

func PrintProperty(name, value string) {
	totalLength := len(name) + len(value)
	if totalLength > 60 {
		fmt.Println(chalk.Yellow, name, "\n", chalk.White, value, chalk.Reset)
	} else {
		spaces := strings.Repeat(" ", 60-totalLength)
		fmt.Println(chalk.Yellow, name, spaces, chalk.White, value, chalk.Reset)
	}
}

func PrintJson(contentBytes []byte) {
	var (
		out     bytes.Buffer
		content string
		err     error
	)
	err = json.Indent(&out, contentBytes, "", "  ")
	if err == nil {
		content = string(out.Bytes())
	} else {
		println(err.Error())
	}

	re := regexp.MustCompile(`([\[\]\{\}]{1})`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Green, chalk.Reset))

	re = regexp.MustCompile(`:\s*\"(.*)\"`)
	content = re.ReplaceAllString(content, fmt.Sprintf(": \"%s$1%s\"", chalk.Blue, chalk.Reset))

	re = regexp.MustCompile(`(\:\s*[true|false]+\s*[,\n])`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Magenta, chalk.Reset))

	re = regexp.MustCompile(`(\:\s*[0-9]+\s*[,\n])`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Yellow, chalk.Reset))

	re = regexp.MustCompile(`(\:\s*null\s*[,\n])`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Red, chalk.Reset))

	fmt.Println(content)
}
