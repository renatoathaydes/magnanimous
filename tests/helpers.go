package tests

import (
	"bufio"
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
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
		diff := shortDiff(actual, expected)
		t.Errorf("[%d] Expected '%s' but was '%s'.\nShort diff:\n%s",
			index, expected, actual, diff)
	}
}

func shortDiff(actual, expected string) string {
	exptdLeft := expected
	eIdx := 0
	for aIdx, a := range actual {
		if len(exptdLeft) > 0 {
			e, size := utf8.DecodeRuneInString(exptdLeft)
			exptdLeft = exptdLeft[size:]
			eIdx += size
			if a != e {
				shortA := actual[uint(math.Max(float64(aIdx-6), 0)):aIdx] + "(" + string(a) + ")"
				shortE := expected[uint(math.Max(float64(eIdx-6), 0)):eIdx] + "(" + string(e) + ")"
				return fmt.Sprintf("at index %d\n%s\n%s", aIdx, shortA, shortE)
			}
		} else {
			return fmt.Sprintf("actual longer than expected, max_index=%d", aIdx)
		}
	}
	return fmt.Sprintf("No differences found")
}

func CreateTempFiles(files map[string]string) (mg.WebFilesMap, string) {
	dir, err := ioutil.TempDir("", "for_test")
	check(err)
	fmt.Printf("Temp dir at %s\n", dir)

	// just create the directory structure with empty files, contents are not required
	filesMap := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	for name, content := range files {
		err = os.MkdirAll(filepath.Join(dir, filepath.Dir(name)), 0770)
		check(err)
		file, err := os.Create(filepath.Join(dir, name))
		check(err)
		_, err = file.Write([]byte(content))
		check(err)
		fileReader := bufio.NewReader(strings.NewReader(content))
		pf, err := mg.ProcessReader(fileReader, name, len(content), nil, time.Now())
		check(err)
		filesMap.WebFiles[filepath.Join(dir, name)] = mg.WebFile{
			Processed:   pf,
			Name:        filepath.Base(name),
			NonWritable: strings.HasPrefix(filepath.Base(name), "_"),
		}
	}

	return filesMap, dir
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkParsing(t *testing.T,
	m mg.WebFilesMap, pf *mg.ProcessedFile,
	expectedCtx map[string]interface{}, expectedContents []string) {

	contents := pf.GetContents()

	if len(contents) != len(expectedContents) {
		t.Fatalf("Expected %d content parts but got %d: %v",
			len(expectedContents), len(contents), contents)
	}

	ctx := mg.NewContext()
	stack := mg.NewContextStack(ctx)

	for i, c := range contents {
		var result strings.Builder
		err := c.Write(&result, m, stack)
		check(err)

		if result.String() != expectedContents[i] {
			t.Errorf("Unexpected Content[%d]\nExpected: '%s'\nActual  : '%s'",
				i, expectedContents[i], result.String())
		}
	}

	if len(expectedCtx) == 0 {
		if !ctx.IsEmpty() {
			t.Errorf("Expected empty context.\n"+
				"Actual Context: %v", ctx)
		}
	} else if !isContextEqual(ctx, expectedCtx) {
		t.Errorf(
			"Expected Context: %v\n"+
				"Actual Context: %v", expectedCtx, ctx)
	}
}

func isContextEqual(context mg.Context, expectedContents map[string]interface{}) bool {
	for key, expectedValue := range expectedContents {
		actualValue, ok := context.Get(key)
		if !ok || !reflect.DeepEqual(expectedValue, actualValue) {
			return false
		}
	}
	return true
}

func checkContents(t *testing.T,
	m mg.WebFilesMap, pf *mg.ProcessedFile,
	expectedContent string) {

	ctx := mg.NewContext()
	stack := mg.NewContextStack(ctx)

	content, err := pf.Bytes(m, stack)

	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expectedContent {
		t.Errorf("Unexpected content. Expected:\n'%s'\nActual:\n'%s'", expectedContent, content)
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

	err = mag.WriteTo(dir, webFiles)
	if err != nil {
		t.Fatal(err)
	}

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
