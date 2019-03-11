package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestIncludeFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("ABCDEF"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ include /processed/hi.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, *resolver.Files, otherProcessed, expectedCtx, []string{
		"OUTER\n",
		"ABCDEF",
		"\nEND"})
}

func TestIncludeEvalFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("ABCDEF"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("{{define p `/processed`}}" +
		"OUTER\n{{ include eval p + `/hi.txt` }}\nEND"))

	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["p"] = "/processed"

	checkParsing(t, *resolver.Files, otherProcessed, expectedCtx, []string{
		"",
		"OUTER\n",
		"ABCDEF",
		"\nEND"})
}

func TestIncludeImplicitEvalFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("ABCDEF"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("{{define p `/processed`}}" +
		"OUTER\n{{ include \"/processed\" + `/hi.txt` }}\nEND"))

	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["p"] = "/processed"

	checkParsing(t, *resolver.Files, otherProcessed, expectedCtx, []string{
		"",
		"OUTER\n",
		"ABCDEF",
		"\nEND"})
}

func TestIncludeFileNested(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("A1"))
	processed1, err := mg.ProcessReader(r, "source/a1.txt", 6, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("A2\n{{include /a1.txt}}"))
	processed2, err := mg.ProcessReader(r, "source/processed/a2.txt", 6, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	// include relative path
	r = bufio.NewReader(strings.NewReader("A3\n{{ include a2.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/a1.txt"] = mg.WebFile{Processed: processed1}
	files["source/processed/a2.txt"] = mg.WebFile{Processed: processed2}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, *resolver.Files, otherProcessed, expectedCtx, []string{
		"A3\n",
		"A2\nA1",
		"\nEND"})
}
