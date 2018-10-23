package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestProcessIncludeMissingCloseBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("## hello {{ include /example.md "))
	_, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	shouldHaveError(t, err, mg.ParseError,
		"(source/processed/hi.md:1:33) instruction started at (1:10) was not properly closed with '}}'")
}

func TestProcessIncludeMissingCloseBracketsAfterGoodInstructions(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("# Doc\n" +
		"## Section 1\n\n" +
		"hello {{ abc name }}\n" +
		"something {{ include abc\n" +
		"footer\n"))
	_, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	shouldHaveError(t, err, mg.ParseError,
		"(source/processed/hi.md:7:1) instruction started at (5:11) was not properly closed with '}}'")
}

func TestInclusionIndirectCycleError(t *testing.T) {
	files := make(mg.WebFilesMap)
	resolver := mg.DefaultFileResolver{BasePath: "", Files: files}

	r := bufio.NewReader(strings.NewReader("A = {{ include /processed/other.txt }}"))
	processed, err := mg.ProcessReader(r, "/processed/hi.txt", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("{{ include /processed/hi.txt }}"))
	otherProcessed, otherErr := mg.ProcessReader(r, "/processed/other.txt", 11, &resolver)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	dir, dirErr := ioutil.TempDir("", "TestInclusionIndirectCycleError")

	if dirErr != nil {
		t.Fatal(dirErr)
	}

	defer os.RemoveAll(dir)

	magErr := mg.WriteTo(dir, files)

	shouldHaveError(t, magErr, mg.InclusionCycleError, "Cycle detected! Inclusion of "+
		"/processed/other.txt at /processed/hi.txt:1:5 "+
		"comes back into itself via [/processed/other.txt:1:1 -> /processed/hi.txt:1:5]",
		"Cycle detected! Inclusion of "+
			"/processed/hi.txt at /processed/other.txt:1:1 "+
			"comes back into itself via [/processed/hi.txt:1:5 -> /processed/other.txt:1:1]")
}

func TestInclusionSelfCycleError(t *testing.T) {
	files := make(mg.WebFilesMap)
	resolver := mg.DefaultFileResolver{BasePath: "", Files: files}

	r := bufio.NewReader(strings.NewReader("A = {{ include hi.txt }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}

	dir, dirErr := ioutil.TempDir("", "TestInclusionSelfCycleError")

	if dirErr != nil {
		t.Fatal(dirErr)
	}

	defer os.RemoveAll(dir)

	magErr := mg.WriteTo(dir, files)

	shouldHaveError(t, magErr, mg.InclusionCycleError, "Cycle detected! Inclusion of "+
		"hi.txt at source/processed/hi.txt:1:5 "+
		"comes back into itself via [source/processed/hi.txt:1:5]")
}
