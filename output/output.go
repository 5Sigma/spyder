package output

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ttacon/chalk"
	"os"
	"regexp"
	"strings"
)

var PromptStream *os.File

var SupressOutput = false

func PrintResult(desc string, err error) {
	resultColor := chalk.Red
	resultText := "FAIL"
	if err == nil {
		resultColor = chalk.Green
		resultText = "OK"
	}
	desc += strings.Repeat(" ", 50-len(desc))
	Println(chalk.White, desc, chalk.Yellow, "[", resultColor, resultText,
		chalk.Yellow, "]", chalk.Reset)
}

func PrintError(err error) {
	Println(chalk.Red, err.Error(), chalk.Reset)
}

func PrintFatal(err error) {
	PrintError(err)
	os.Exit(1)
}

func PrintProperty(name, value string) {
	totalLength := len(name) + len(value)
	if totalLength > 60 {
		Println(chalk.Yellow, name, "\n", chalk.White, value, chalk.Reset)
	} else {
		spaces := strings.Repeat(" ", 60-totalLength)
		Println(chalk.Yellow, name, spaces, chalk.White, value, chalk.Reset)
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
		PrintError(err)
	}

	re := regexp.MustCompile(`([\[\]\{\}]{1})`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Green, chalk.Reset))

	// String values
	re = regexp.MustCompile(`(\s*?\")([^:]*?)(\"\s*?,?\n)`)
	content = re.ReplaceAllString(content, fmt.Sprintf("$1%s$2%s$3", chalk.Blue, chalk.Reset))

	re = regexp.MustCompile(`(\:\s*[true|false]+\s*[,\n])`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Magenta, chalk.Reset))

	re = regexp.MustCompile(`(\:\s*[0-9]+\s*[,\n])`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Yellow, chalk.Reset))

	re = regexp.MustCompile(`(\:\s*null\s*[,\n])`)
	content = re.ReplaceAllString(content, fmt.Sprintf("%s$1%s", chalk.Red, chalk.Reset))

	Println(content)
}

func Prompt(name, defaultValue string) string {
	if PromptStream == nil {
		PromptStream = os.Stdin
	}
	reader := bufio.NewReader(PromptStream)
	Printf("%s%s [%s]: %s", chalk.Yellow, name, defaultValue, chalk.Reset)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func Println(str ...interface{}) {
	if SupressOutput == false {
		fmt.Println(str...)
	}
}

func Printf(str string, args ...interface{}) {
	if SupressOutput == false {
		fmt.Printf(str, args...)
	}
}
