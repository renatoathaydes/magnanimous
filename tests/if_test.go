package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestIfTrue(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if true }}" +
		"YES" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\nYES")
}

func TestIfTrueExpr(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if 2 + 2 == 4 }}" +
		"YES" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\nYES")
}

func TestIfFalse(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if false }}" +
		"NO" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\n")
}

func TestIfNil(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if nil }}" +
		"NO" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\n")
}

func TestIfFalseExpr(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if 2 > 100 }}" +
		"NO" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\n")
}

func TestIfNegatedFalseExpr(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if !(2 + 2 > 10) }}" +
		"NO" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\nNO")
}

func TestIfNegatedTrueExpr(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if !(2 + 2 < 10) }}" +
		"NO" +
		"{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\n")
}

func TestIfNonBooleanCondition(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("IF:\n" +
		"{{ if 10 }}INT{{ end }}" +
		"{{ if \"hi\" }}STRING{{ end }}" +
		"{{ if 2 * 2 }}MULT{{ end }}" +
		"{{ if 0 }}ZERO{{ end }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "IF:\n")
}
