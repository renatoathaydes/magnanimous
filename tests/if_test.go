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

func TestIfScope(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ define x 10 }}\n" +
			"Before IF, X = {{ eval x }}\n" +
			"{{ if 2 + 2 == 4 }}\n" +
			"  Inside IF, X = {{ define x 20 }}{{ eval x }}\n" +
			"{{ define y 2 }}\n" +
			"  Inside IF, Y = {{ eval y }}\n" +
			"{{ end }}\n" +
			"After IF, X = {{ eval x }} and Y = {{ eval y }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkContents(t, emptyFilesMap, processed, "\nBefore IF, X = 10\n\n"+
		"  Inside IF, X = 20\n\n"+
		"  Inside IF, Y = 2\n\n"+
		"After IF, X = 10 and Y = ")
}
