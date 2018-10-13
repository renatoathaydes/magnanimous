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

	if err == nil {
		t.Fatal("No error occurred!")
	}

	if err.Code != mg.ParseError {
		t.Errorf("Expected ParseError but got %s\n", err.Code)
	}

	if err.Error() != "(source/processed/hi.md:1:32) instruction started at (1:10) was not properly closed with '}}'" {
		t.Errorf("Unexpected Error(): %s\n", err.Error())
	}
}

func TestProcessIncludeMissingCloseBracketsAfterGoodInstructions(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("# Doc\n" +
		"## Section 1\n\n" +
		"hello {{ eval name }}\n" +
		"something {{ include abc\n" +
		"footer\n"))
	_, _, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err == nil {
		t.Fatal("No error occurred!")
	}

	if err.Code != mg.ParseError {
		t.Errorf("Expected ParseError but got %s\n", err.Code)
	}

	if err.Error() != "(source/processed/hi.md:7:1) instruction started at (5:12) was not properly closed with '}}'" {
		t.Errorf("Unexpected Error(): %s\n", err.Error())
	}
}
