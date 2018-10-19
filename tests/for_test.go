package tests

import (
	"bufio"
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestForArray(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for var (1,2,3, 42) }}\n" +
		"Number {{ eval var }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed,
		"Loop Sample:\n\n"+
			"Number 1\n\n"+
			"Number 2\n\n"+
			"Number 3\n\n"+
			"Number 42\n")
}

func TestForArrayInMarkDown(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ for section (\"Home\", \"About\") }}\n" +
			"## {{ eval section }}\nSomething something{{ end }}\n" +
			"END"))
	processed, err := mg.ProcessReader(r, "source/processed/array.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed,
		"<h2>Home</h2>\n\n"+
			"<p>Something something</p>\n\n"+
			"<h2>About</h2>\n\n"+
			"<p>Something something</p>\n"+
			"<p>END</p>\n")
}

func TestForFiles(t *testing.T) {
	// replace temporarily Magnanimous' FileResolver with the test one
	defaultFileResolver := mg.DefaultFileResolver
	defer func() { mg.DefaultFileResolver = defaultFileResolver }()
	mg.DefaultFileResolver = &TestFileResolver

	// create a bunch of files for testing
	files := CreateTempFiles()
	defer DeleteTempFiles()

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path /processed/examples }}\n" +
		"Title {{ eval path }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11)

	if err != nil {
		t.Fatal(err)
	}

	// FIXME expression can't evaluate Map properties, so we're just putting the Map itself in the template
	checkContents(t, files, processed,
		"Loop Sample:\n\n"+
			"Title map[title:File 1]\n\n"+
			"Title map[title:Second File]\n")
}

var TestFileResolver = resolver{}

func CreateTempFiles() mg.WebFilesMap {
	dir, err := ioutil.TempDir("", "for_test")
	check(err)
	TestFileResolver.tempDir = dir
	fmt.Printf("Temp dir at %s\n", dir)

	// just create the directory structure with empty files, contents are not required
	err = os.MkdirAll(filepath.Join(dir, "processed/examples"), 0770)
	check(err)
	_, err = os.Create(filepath.Join(dir, "processed/examples/f1.txt"))
	check(err)
	_, err = os.Create(filepath.Join(dir, "processed/examples/f2.txt"))
	check(err)

	files := mg.WebFilesMap{}

	files["/processed/examples/f1.txt"] = mg.WebFile{Processed: &mg.ProcessedFile{}}
	files["/processed/examples/f1.txt"].Processed.Context()["title"] = "File 1"

	files["/processed/examples/f2.txt"] = mg.WebFile{Processed: &mg.ProcessedFile{}}
	files["/processed/examples/f2.txt"].Processed.Context()["title"] = "Second File"

	return files
}

func DeleteTempFiles() {
	if len(TestFileResolver.tempDir) > 0 {
		os.RemoveAll(TestFileResolver.tempDir)
		TestFileResolver.tempDir = ""
	}
}

type resolver struct {
	tempDir string
}

func (r resolver) FilesIn(dir string, from mg.Location) (dirPath string, f []os.FileInfo, e error) {
	dirPath = r.Resolve(dir, from)
	fmt.Printf("Resolved %s for %s\n", dirPath, dir)
	f, e = ioutil.ReadDir(filepath.Join(r.tempDir, dirPath))
	return
}

func (resolver) Resolve(path string, from mg.Location) string {
	return path
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
