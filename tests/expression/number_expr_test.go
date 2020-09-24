package expression

import (
	"reflect"
	"testing"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

func TestIntExpr(t *testing.T) {
	e, err := expression.ParseExpr(`100`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(&e, nil)

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

	v, err := expression.EvalExpr(&e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(52) {
		t.Errorf("Expected '52' but got '%v'", v)
	}
}

func TestIntSubtractionExpr(t *testing.T) {
	v, err := expression.Eval(`5 - 2`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(3) {
		t.Errorf("Expected '3' but got '%v'", v)
	}
}

func TestIntMultiplicationExpr(t *testing.T) {
	v, err := expression.Eval(`3*6`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(18) {
		t.Errorf("Expected '18' but got '%v'", v)
	}
}

func TestIntDivisionExpr(t *testing.T) {
	v, err := expression.Eval(`10 / 5`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(2) {
		t.Errorf("Expected '2' but got '%v'", v)
	}
}

func TestIntRemainderExpr(t *testing.T) {
	v, err := expression.Eval(`10 % 4`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(2) {
		t.Errorf("Expected '2' but got '%v'", v)
	}
}

func TestComplexArithmeticExpr(t *testing.T) {
	v, err := expression.Eval(`(2 + 4) * 5 / 10`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(3) {
		t.Errorf("Expected '3' but got '%v'", v)
	}
}

func TestComplexArithmeticExpr2(t *testing.T) {
	v, err := expression.Eval(`( 5 * 2 * 10 / 4 ) + ( 2 * 10 / 4 )`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(30) {
		t.Errorf("Expected '30' but got '%v'", v)
	}
}

func TestIntArrayExpr(t *testing.T) {
	v, err := expression.Eval(`[]interface{}{1,3,5}`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if !reflect.DeepEqual(v, []interface{}{float64(1), float64(3), float64(5)}) {
		t.Errorf("Expected '[1,3,5]' but got '%v'", v)
	}
}

func TestIntEqualityExpr(t *testing.T) {
	v, err := expression.Eval(`2 == 2`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestIntInequalityExpr(t *testing.T) {
	v, err := expression.Eval(`2 != 2`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != false {
		t.Errorf("Expected 'false' but got '%v'", v)
	}
}

func TestIntGreaterThanExpr(t *testing.T) {
	v, err := expression.Eval(`5 > 1`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestIntGreaterThanOrEqualExpr(t *testing.T) {
	v, err := expression.Eval(`5 >= 1`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestIntLessThanExpr(t *testing.T) {
	v, err := expression.Eval(`10 < 100`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestIntLessThanOrEqualExpr(t *testing.T) {
	v, err := expression.Eval(`10 <= 10`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestIntMinusSign(t *testing.T) {
	v, err := expression.Eval(`-19`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(-19) {
		t.Errorf("Expected '-19' but got '%v'", v)
	}
}

func TestIntPlusSign(t *testing.T) {
	v, err := expression.Eval(`+53`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != float64(53) {
		t.Errorf("Expected '+53' but got '%v'", v)
	}
}
