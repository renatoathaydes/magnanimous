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

func TestIntExpr(t *testing.T) {
	e, err := expression.ParseExpr(`100`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(100) {
		t.Errorf("Expected '100' but got '%v'", v)
	}
}

func TestIntAdditionExpr(t *testing.T) {
	e, err := expression.ParseExpr(`42 + 10`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(52) {
		t.Errorf("Expected '52' but got '%v'", v)
	}
}
