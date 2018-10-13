package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestProcessIncludeMissingCloseBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("## hello {{ include /example.md "))
	_, _, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	shouldHaveError(t, err, mg.ParseError,
		"(source/processed/hi.md:1:33) instruction started at (1:10) was not properly closed with '}}'")
}

func TestProcessIncludeMissingCloseBracketsAfterGoodInstructions(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("# Doc\n" +
		"## Section 1\n\n" +
		"hello {{ abc name }}\n" +
		"something {{ include abc\n" +
		"footer\n"))
	_, _, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	shouldHaveError(t, err, mg.ParseError,
		"(source/processed/hi.md:7:1) instruction started at (5:11) was not properly closed with '}}'")
}

func shouldHaveError(t *testing.T, err *mg.MagnanimousError, code mg.ErrorCode, message string) {
	if err == nil {
		t.Fatal("No error occurred!")
	}
	if err.Code != code {
		t.Errorf("Expected %s but got %s\n", code, err.Code)
	}
	if err.Error() != message {
		t.Errorf("Unexpected error message.\nExpected: %s\nActual  : %s", message, err.Error())
	}
}
