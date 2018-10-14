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
	_, _, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	shouldHaveError(t, err, mg.ParseError,
		"(source/processed/hi.md:1:33) instruction started at (1:10) was not properly closed with '}}'")
}

func TestProcessIncludeMissingCloseBracketsAfterGoodInstructions(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("# Doc\n" +
		"## Section 1\n\n" +
		"hello {{ abc name }}\n" +
		"something {{ include abc\n" +
		"footer\n"))
	_, _, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	shouldHaveError(t, err, mg.ParseError,
		"(source/processed/hi.md:7:1) instruction started at (5:11) was not properly closed with '}}'")
}

func TestInclusionIndirectCycleError(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("A = {{ include /processed/other.txt }}"))
	_, processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("{{ include /processed/hi.txt }}"))
	otherCtx, otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files := mg.WebFilesMap{}
	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed, Context: emptyContext}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed, Context: otherCtx}

	dir, dirErr := ioutil.TempDir("", "TestInclusionIndirectCycleError")

	if dirErr != nil {
		t.Fatal(dirErr)
	}

	defer os.RemoveAll(dir)

	magErr := mg.WriteTo(dir, files)

	shouldHaveError(t, magErr, mg.InclusionCycleError, "Cycle detected! Inclusion of "+
		"source/processed/hi.txt at source/processed/other.txt:1:1 "+
		"comes back into itself via [source/processed/hi.txt:1:5 -> source/processed/other.txt:1:1]")
}

func TestInclusionSelfCycleError(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("A = {{ include hi.txt }}"))
	_, processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11)

	if err != nil {
		t.Fatal(err)
	}

	files := mg.WebFilesMap{}
	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed, Context: emptyContext}

	dir, dirErr := ioutil.TempDir("", "TestInclusionSelfCycleError")

	if dirErr != nil {
		t.Fatal(dirErr)
	}

	defer os.RemoveAll(dir)

	magErr := mg.WriteTo(dir, files)

	shouldHaveError(t, magErr, mg.InclusionCycleError, "Cycle detected! Inclusion of "+
		"source/processed/hi.txt at source/processed/hi.txt:1:5 "+
		"comes back into itself via [source/processed/hi.txt:1:5]")
}

func shouldHaveError(t *testing.T, err *mg.MagnanimousError, code mg.ErrorCode, message string) {
	if err == nil {
		t.Fatal("No error occurred!")
	}
	if err.Code != code {
		t.Errorf("Expected %s but got %s\n", code, err.Code)
	}
	if err.Error() != message {
		t.Errorf("Unexpected error message.\nExpected: %s\nActual  : %s", message, err.Error())
	}
}
