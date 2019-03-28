package expression

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"testing"
	"time"
)

func TestDateExpr(t *testing.T) {
	expected := "2017-02-03"
	checkDate(t, fmt.Sprintf(`date["%s"]`, expected), expected, "2006-01-02")
}

func TestDateExprMid(t *testing.T) {
	expected := "2017-02-03T10:40"
	checkDate(t, fmt.Sprintf(`date["%s"]`, expected), expected, "2006-01-02T15:04")
}

func TestDateExprLong(t *testing.T) {
	expected := "2017-02-03T10:40:32"
	checkDate(t, fmt.Sprintf(`date["%s"]`, expected), expected, "2006-01-02T15:04:05")
}

// Go format: Mon Jan 2 15:04:05 -0700 MST 2006

func checkDate(t *testing.T, dateExpr, expected, format string) {
	e, err := expression.ParseExpr(dateExpr)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	expectedTime, err := time.Parse(format, expected)
	if err != nil {
		panic(err)
	}

	if v != expectedTime {
		t.Errorf("Expected '%v' but got '%v'", expectedTime, v)
	}
}
