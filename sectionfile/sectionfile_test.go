package sectionfile

import (
	"testing"
)

func TestDefaultSection(t *testing.T) {
	contents := `file contents
without any kind of marker but has an @sign in it.`
	sf := LoadString("", contents)
	if len(sf.Sections) == 0 {
		t.Error("No default section loaded")
	}
	if len(sf.Sections) > 1 {
		t.Error("Too many sections loaded")
	}
	if sf.Sections["endpoint"] != contents {
		t.Errorf("contents dont match\n%s\n---\n%s", sf.Sections["endpoint"], contents)
	}
}

func TestSingleSection(t *testing.T) {
	contents := `
		@testname
		file contents
		without any kind of marker but has an @sign in it.
	`
	sf := LoadString("", contents)

	if len(sf.Sections) > 1 {
		t.Error("Too many sections loaded")
	}
	if _, ok := sf.Sections["testname"]; !ok {
		t.Error("Section not named correctly")
	}

}

func TestMultipleSections(t *testing.T) {
	contents := `
@testname
file contents
without any kind of marker but has an @sign in it.
@section2
heres something
@section3
third section`
	sf := LoadString("", contents)

	if len(sf.Sections) != 3 {
		t.Errorf("%d Sections loaded, should be %d", len(sf.Sections), 3)
	}

	if _, ok := sf.Sections["testname"]; !ok {
		t.Error("Section testname not named correctly")
	}

	if _, ok := sf.Sections["section2"]; !ok {
		t.Error("Section section2 not named correctly")
	}

	if _, ok := sf.Sections["section3"]; !ok {
		t.Error("Section section3 not named correctly")
	}

	if sf.Contents("section2") != "heres something\n" {
		t.Errorf("Contents not correct:\n%s", sf.Contents("section2"))
	}

}
