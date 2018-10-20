package expression

import (
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"testing"
)

func TestStringExpr(t *testing.T) {
	e, err := expression.ParseExpr(`"hello"`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != `hello` {
		t.Errorf("Expected 'hello' but got '%v'", v)
	}
}

func TestSingleQuoteStringExpr(t *testing.T) {
	e, err := expression.ParseExpr("`hi\nthere`")

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "hi\nthere" {
		t.Errorf("Expected 'hi\nthere' but got '%v'", v)
	}
}

func TestStringAdditionExpr(t *testing.T) {
	e, err := expression.ParseExpr(`"ab" + "cd"`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "abcd" {
		t.Errorf("Expected 'abcd' but got '%v'", v)
	}
}

func TestStringEqualityExpr(t *testing.T) {
	v, err := expression.Eval(`"Boo" == "Boo"`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestStringInequalityExpr(t *testing.T) {
	v, err := expression.Eval(`"FOO" != "BAR"`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestStringGreaterThanExpr(t *testing.T) {
	v, err := expression.Eval(`"z" > "a"`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestStringGreaterThanOrEqualExpr(t *testing.T) {
	v, err := expression.Eval(`"z" >= "a"`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestStringLessThanExpr(t *testing.T) {
	v, err := expression.Eval(`"a" < "z"`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestStringLessThanOrEqualExpr(t *testing.T) {
	v, err := expression.Eval(`"a" <= "a"`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}
