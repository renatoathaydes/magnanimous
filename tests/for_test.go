package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
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
