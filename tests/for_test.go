package tests

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/renatoathaydes/magnanimous/mg"
)

func TestForArray(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for v [1,2,3, 42] }}\n" +
		"Number {{ eval v }}\\\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Number 1\n"+
			"Number 2\n"+
			"Number 3\n"+
			"Number 42")
}

func TestForDefineArray(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ define a [1,2,3, 42] }}" +
		"{{for v eval a}}" +
		"Number {{ eval v }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n"+
			"Number 1\n"+
			"Number 2\n"+
			"Number 3\n"+
			"Number 42\n")
}

func TestForNestedArray(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Nested Sample:\n" +
		"{{ for i [1,2] }}" +
		"{{ for j [6,5] }}" +
		"Numbers {{eval i}} {{eval j}}\n" +
		"{{ end }}" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Nested Sample:\n"+
			"Numbers 1 6\n"+
			"Numbers 1 5\n"+
			"Numbers 2 6\n"+
			"Numbers 2 5\n")
}

func TestForArrayWithExpressions(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:" +
		"{{ for x [1 + 1, 2 + 2] }}\n" +
		"X is {{ eval x }}" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Numbers:\n"+
			"X is 2\n"+
			"X is 4")
}

func TestForArraySorted(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:\n" +
		"{{ for x ( sort ) [10, 2, 4, 1, 2, 5] }}" +
		"{{ eval x }} " +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Numbers:\n"+
			"1 2 2 4 5 10 ")
}

func TestForArraySortedBy(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:\n" +
		"{{ for x (sortBy _) [10, 2, 4, 1, 2, 5] }}" +
		"{{ eval x }} " +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Numbers:\n"+
			"1 2 2 4 5 10 ")
}

func TestForArrayReverse(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:\n" +
		"{{ for x ( reverse ) [10, 2, 4, 1, 2, 5] }}" +
		"{{ eval x }} " +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Numbers:\n"+
			"5 2 1 4 2 10 ")
}

func TestForArraySortReverse(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:\n" +
		"{{ for x ( sort reverse ) [10, 2, 4, 1, 2, 5] }}" +
		"{{ eval x }} " +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Numbers:\n"+
			"10 5 4 2 2 1 ")
}

func TestForArrayLimit(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for v (limit 3) [ 1 , 2, 3, 4 , 5 , 6 ] }}\n" +
		"Number {{ eval v }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Number 1\n\n"+
			"Number 2\n\n"+
			"Number 3\n")
}

func TestForArraySortByReverseLimit(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers:\n" +
		"{{ for x ( sortBy _ reverse limit 4 ) [10, 2, 4, 1, 2, 5] }}" +
		"{{ eval x }} " +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Numbers:\n"+
			"10 5 4 2 ")
}

func TestForArrayInMarkDown(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ for section [ \"Home\", \"About\" ] }}\n" +
			"## {{ eval section }}\nSomething something{{ end }}\n" +
			"END"))
	processed, err := mg.ProcessReader(r, "source/processed/array.md", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"<h2>Home</h2>\n\n"+
			"<p>Something something</p>\n\n"+
			"<h2>About</h2>\n\n"+
			"<p>Something something\n"+
			"END</p>\n")
}

func TestForFiles(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{ define title \"File 1\" }}",
		"processed/examples/f2.txt": "{{ define title \"Second File\" }}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path /processed/examples }}\n" +
		"Title {{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Title File 1\n\n"+
			"Title Second File\n")
}

func TestForFilesScope(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{ define title \"File 1\" }}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Title is {{ eval title }}\n" +
		"{{ for path /processed/examples }}\n" +
		"  path.title: {{ eval path.title }}\n" +
		"  title: {{ eval title }}\n" +
		"{{ end }}\n" +
		"Title is {{ eval title }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Title is \n\n"+
			"  path.title: File 1\n"+
			"  title: \n\n"+
			"Title is ")
}

func TestForFilesWithUnwritableFiles(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{ define title \"File 1\" }}",
		"processed/examples/_a.txt": "{{ define title \"?\" }}",
		"processed/examples/f2.txt": "{{ define title \"Second File\" }}",
		"processed/examples/_b.txt": "{{ define title \"?\" }}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path /processed/examples }}\n" +
		"Title {{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Title File 1\n\n"+
			"Title Second File\n")
}

func TestForFilesReverse(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/a.txt": "{{ define title \"A\"}}",
		"processed/b.txt": "{{ define title \"B\"}}",
		"processed/g.txt": "{{ define title \"G\"}}",
		"processed/f.txt": "{{ define title \"F\"}}",
		"processed/c.txt": "{{ define title \"C\"}}",
		"processed/e.txt": "{{ define title \"E\"}}",
		"processed/d.txt": "{{ define title \"D\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path ( reverse ) /processed/ }}\n" +
		"Title {{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Title G\n\n"+
			"Title F\n\n"+
			"Title E\n\n"+
			"Title D\n\n"+
			"Title C\n\n"+
			"Title B\n\n"+
			"Title A\n")
}

func TestForFilesLimit(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/a.txt": "{{define title \"A\"}}",
		"processed/b.txt": "{{define title \"B\"}}",
		"processed/g.txt": "{{define title \"G\"}}",
		"processed/f.txt": "{{define title \"F\"}}",
		"processed/c.txt": "{{define title \"C\"}}",
		"processed/e.txt": "{{define title \"E\"}}",
		"processed/d.txt": "{{define title \"D\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path ( limit 5 ) /processed/ }}\n" +
		"Title {{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Title A\n\n"+
			"Title B\n\n"+
			"Title C\n\n"+
			"Title D\n\n"+
			"Title E\n")
}

func TestForFilesLimitTooMany(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/a.txt": "{{define title \"A\"}}",
		"processed/b.txt": "{{define title \"B\"}}",
		"processed/c.txt": "{{define title \"C\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path ( limit 56 ) /processed/ }}\n" +
		"Title {{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Title A\n\n"+
			"Title B\n\n"+
			"Title C\n")
}

func TestForFilesSortBy(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{define title \"Some file\"}}",
		"processed/examples/f2.txt": "{{define title \"Other file\"}}",
		"processed/examples/f3.txt": "{{define title \"A file\"}}",
		"processed/examples/f4.txt": "{{define title \"Z file\"}}",
		"processed/examples/f5.txt": "{{define title \"Final\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path (sortBy title) /processed/examples }}\n" +
		"{{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"A file\n\n"+
			"Final\n\n"+
			"Other file\n\n"+
			"Some file\n\n"+
			"Z file\n")
}

func TestForFilesSortByReverse(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{define title \"Some file\"}}",
		"processed/examples/f2.txt": "{{define title \"Other file\"}}",
		"processed/examples/f3.txt": "{{define title \"A file\"}}",
		"processed/examples/f4.txt": "{{define title \"Z file\"}}",
		"processed/examples/f5.txt": "{{define title \"Final\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path (sortBy title reverse) /processed/examples }}\n" +
		"{{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Z file\n\n"+
			"Some file\n\n"+
			"Other file\n\n"+
			"Final\n\n"+
			"A file\n")
}

func TestForFilesReverseSortBy(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{define title \"Some file\"}}",
		"processed/examples/f2.txt": "{{define title \"Other file\"}}",
		"processed/examples/f3.txt": "{{define title \"A file\"}}",
		"processed/examples/f4.txt": "{{define title \"Z file\"}}",
		"processed/examples/f5.txt": "{{define title \"Final\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path (reverse   sortBy title) /processed/examples }}\n" +
		"{{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"A file\n\n"+
			"Final\n\n"+
			"Other file\n\n"+
			"Some file\n\n"+
			"Z file\n")
}

func TestForFilesLimitSortByReverse(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{define title \"Some file\"}}",
		"processed/examples/f2.txt": "{{define title \"Other file\"}}",
		"processed/examples/f3.txt": "{{define title \"A file\"}}",
		"processed/examples/f4.txt": "{{define title \"Z file\"}}",
		"processed/examples/f5.txt": "{{define title \"Final\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ for path ( limit 3 sortBy title reverse ) /processed/examples }}\n" +
		"{{ eval path.title }}\n" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Some file\n\n"+
			"Other file\n\n"+
			"A file\n")
}

func TestForFilesEval(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{define title \"Some file\"}}",
		"processed/examples/f2.txt": "{{define title \"Other file\"}}",
		"processed/examples/f3.txt": "{{define title \"A file\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ define dir `examples` }}" +
		"{{ for path eval `/processed/` + dir }}\n" +
		"{{ eval path.title }}\n" +
		"{{ end }}"))

	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"Some file\n\n"+
			"Other file\n\n"+
			"A file\n")
}

func TestForFilesEvalSortBy(t *testing.T) {

	// create a bunch of files for testing
	files, dir := CreateTempFiles(map[string]string{
		"processed/examples/f1.txt": "{{define title \"Some file\"}}",
		"processed/examples/f2.txt": "{{define title \"Other file\"}}",
		"processed/examples/f3.txt": "{{define title \"A file\"}}",
	})
	defer os.RemoveAll(dir)

	resolver := mg.DefaultFileResolver{BasePath: dir, Files: &files}

	r := bufio.NewReader(strings.NewReader("Loop Sample:\n" +
		"{{ define top `/processed` }}" +
		"{{ define dir `/examples` }}" + // intentional double-slash below
		"{{ for path (sortBy title) eval top + `/` + dir }}\n" +
		"{{ eval path.title }}\n" +
		"{{ end }}"))

	processed, err := mg.ProcessReader(r, filepath.Join(dir, "processed/hi.txt"), dir, 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, processed,
		"Loop Sample:\n\n"+
			"A file\n\n"+
			"Other file\n\n"+
			"Some file\n")
}
