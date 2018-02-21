package sectionfile

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// SectionFile a representation of the file with its content broke out into
// sections. Sections are denoted by a line begining with a '@' and then a
// section name
type SectionFile struct {
	Filename string
	Sections map[string]string
}

var defaultName = "endpoint"

var headerRx, _ = regexp.Compile(`\A\s*@([A-Za-z0-9]+)\z`)

// Load - Loads from a file on the disk.
func Load(filename string) (*SectionFile, error) {
	var (
		fileBytes []byte
		err       error
	)
	fileBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadString(filename, string(fileBytes)), nil
}

func LoadString(filename, fileContents string) *SectionFile {
	var (
		sectionName    string
		sectionContent string
		lines          []string
		sf             *SectionFile
	)

	sf = &SectionFile{Filename: filename, Sections: map[string]string{}}
	sectionName = defaultName
	sectionContent = ""
	lines = strings.Split(fileContents, "\n")

	for idx, line := range lines {

		if strings.TrimSpace(line) == "" {
			continue
		}

		matches := headerRx.FindStringSubmatch(line)

		if len(matches) > 0 {
			if strings.TrimSpace(sectionContent) != "" {
				sf.Sections[sectionName] = sectionContent
			}
			sectionName = matches[1]
			sectionContent = ""
		} else {
			sectionContent += strings.TrimSpace(line)
			if idx < len(lines)-1 {
				sectionContent += "\n"
			}
		}
	}
	sf.Sections[sectionName] = sectionContent

	return sf
}

// Contents - returns the contents for a section. This method is safe to call
// for nonexisting sections and will return an empty string.
func (sf *SectionFile) Contents(name string) string {
	if _, ok := sf.Sections[name]; ok {
		return sf.Sections[name]
	}
	return ""
}
