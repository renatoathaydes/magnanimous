package tests

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var emptyContext = make(map[string]interface{})
var emptyFilesMap = mg.WebFilesMap{}

func shouldHaveError(t *testing.T, err error, code mg.ErrorCode, messageAlternatives ...string) {
	if err == nil {
		t.Fatal("No error occurred!")
	}
	merr, ok := err.(*mg.MagnanimousError)
	if !ok {
		t.Fatalf("Expected error of type MagnanimousError, but found other type: %v", merr)
	}
	if merr.Code != code {
		t.Errorf("Expected %s but got %s\n", code, merr.Code)
	}
	matchFound := false
	for _, expectedMessage := range messageAlternatives {
		if err.Error() == expectedMessage {
			matchFound = true
			break
		}
	}
	if !matchFound {
		t.Errorf("Unexpected error message. Expected one of:\n" +
			strings.Join(messageAlternatives, "\n    OR\n") + "\n    BUT got:\n" + err.Error())
	}
}

func verifyEqual(index uint16, t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("[%d] Expected '%s' but was '%s'", index, expected, actual)
	}
}

func CreateTempFiles() (mg.WebFilesMap, string) {
	dir, err := ioutil.TempDir("", "for_test")
	check(err)
	fmt.Printf("Temp dir at %s\n", dir)

	// just create the directory structure with empty files, contents are not required
	err = os.MkdirAll(filepath.Join(dir, "processed/examples"), 0770)
	check(err)
	_, err = os.Create(filepath.Join(dir, "processed/examples/f1.txt"))
	check(err)
	_, err = os.Create(filepath.Join(dir, "processed/examples/f2.txt"))
	check(err)

	files := mg.WebFilesMap{}

	files[filepath.Join(dir, "processed/examples/f1.txt")] = mg.WebFile{Processed: &mg.ProcessedFile{}, Name: "f1.txt"}
	files[filepath.Join(dir, "processed/examples/f1.txt")].Processed.Context()["title"] = "File 1"

	files[filepath.Join(dir, "processed/examples/f2.txt")] = mg.WebFile{Processed: &mg.ProcessedFile{}, Name: "f2.txt"}
	files[filepath.Join(dir, "processed/examples/f2.txt")].Processed.Context()["title"] = "Second File"

	return files, dir
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkParsing(t *testing.T,
	ctx map[string]interface{}, m mg.WebFilesMap, pf *mg.ProcessedFile,
	expectedCtx map[string]interface{}, expectedContents []string) {

	if len(pf.Contents) != len(expectedContents) {
		t.Fatalf("Expected %d content parts but got %d: %v",
			len(expectedContents), len(pf.Contents), pf.Contents)
	}

	for i, c := range pf.Contents {
		var result strings.Builder
		c.Write(&result, m, nil)

		if result.String() != expectedContents[i] {
			t.Errorf("Unexpected Content[%d]\nExpected: '%s'\nActual  : '%s'",
				i, expectedContents[i], result.String())
		}
	}

	if len(expectedCtx) == 0 {
		if len(ctx) != 0 {
			t.Errorf("Expected empty context.\n"+
				"Actual Context: %v", ctx)
		}
	} else if !reflect.DeepEqual(ctx, expectedCtx) {
		t.Errorf(
			"Expected Context: %v\n"+
				"Actual Context: %v", expectedCtx, ctx)
	}
}

func checkContents(t *testing.T,
	m mg.WebFilesMap, pf *mg.ProcessedFile,
	expectedContent string) {

	content, err := pf.Bytes(m, nil)

	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expectedContent {
		t.Errorf("Unexpected content. Expected:\n%s\nActual:\n%s", expectedContent, content)
	}

}

func runMg(t *testing.T, project string) string {
	mag := mg.Magnanimous{SourcesDir: project}
	webFiles, err := mag.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", project)
	if err != nil {
		t.Fatal(err)
	}

	mg.WriteTo(dir, webFiles)

	return dir
}

func assertFileContents(t *testing.T, files []string, baseDir, file, expectedContent string) {
	var targetFile *string
	for _, f := range files {
		if f == file {
			targetFile = &f
			break
		}
	}
	if targetFile != nil {
		c, err := ioutil.ReadFile(filepath.Join(baseDir, *targetFile))
		if err == nil {
			verifyEqual(0, t, string(c), expectedContent)
		} else {
			t.Fatalf("Error reading file %s: %v\n", file, err)
		}
	} else {
		t.Fatalf("Could not find file %s in %v\n", file, files)
	}
}

func readAll(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			fPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, fPath)
		}
		return err
	})

	return files, err
}
