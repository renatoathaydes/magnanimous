package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
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

func TestForArrayWithExpressions(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:" +
		"{{ for x [1 + 1, 2 + 2] }}\n" +
		"X is {{ eval x }}" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed,
		"Numbers:\n"+
			"X is 2\n"+
			"X is 4")
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
		"Title {{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, files, processed,
		"Loop Sample:\n\n"+
			"Title File 1\n\n"+
			"Title Second File\n")
}
