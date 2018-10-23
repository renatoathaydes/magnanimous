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
		"{{ for v [1,2,3, 42] }}\n" +
		"Number {{ eval v }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

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
		"{{ for section [ \"Home\", \"About\" ] }}\n" +
			"## {{ eval section }}\nSomething something{{ end }}\n" +
			"END"))
	processed, err := mg.ProcessReader(r, "source/processed/array.md", 11, nil)

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

	// create a bunch of files for testing
	files, dir := CreateTempFiles()
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path /processed/examples }}\n" +
		"Title {{ eval path }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	// FIXME expression can't evaluate Map properties, so we're just putting the Map itself in the template
	checkContents(t, files, processed,
		"Loop Sample:\n\n"+
			"Title map[title:File 1]\n\n"+
			"Title map[title:Second File]\n")
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

	files[filepath.Join(dir, "processed/examples/f1.txt")] = mg.WebFile{Processed: &mg.ProcessedFile{}}
	files[filepath.Join(dir, "processed/examples/f1.txt")].Processed.Context()["title"] = "File 1"

	files[filepath.Join(dir, "processed/examples/f2.txt")] = mg.WebFile{Processed: &mg.ProcessedFile{}}
	files[filepath.Join(dir, "processed/examples/f2.txt")].Processed.Context()["title"] = "Second File"

	return files, dir
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
