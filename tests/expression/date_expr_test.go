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

func TestDateLessThanAndGreaterThan(t *testing.T) {
	examples := [][2]string{
		// lesser date, greater date
		{"1990-01-01", "1990-01-02"},
		{"1990-01-01", "1990-02-01"},
		{"1990-01-01", "1991-01-01"},
		{"2017-02-03T10:40:30", "2017-02-03T10:40:31"},
		{"2017-02-03T10:40:30", "2017-02-03T10:44:30"},
		{"2017-02-03T10:40:30", "2017-02-03T12:40:30"},
	}

	for i, ex := range examples {
		// check x < y is true
		expr := fmt.Sprintf("date[\"%s\"] < date[\"%s\"]", ex[0], ex[1])
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != true {
			t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
		}

		// check y < x is false
		expr = fmt.Sprintf("date[\"%s\"] < date[\"%s\"]", ex[1], ex[0])
		v, err = expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != false {
			t.Errorf("[%d] Expected 'false' but got '%v'", i, v)
		}

		// check x > y is false
		expr = fmt.Sprintf("date[\"%s\"] > date[\"%s\"]", ex[0], ex[1])
		v, err = expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != false {
			t.Errorf("[%d] Expected 'false' but got '%v'", i, v)
		}

		// check y > x is true
		expr = fmt.Sprintf("date[\"%s\"] > date[\"%s\"]", ex[1], ex[0])
		v, err = expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != true {
			t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
		}
	}
}

func TestDateLessOrEqual(t *testing.T) {
	examples := [][2]string{
		// lesser date, greater date
		{"1990-01-01", "1990-01-01"},
		{"1990-01-01", "1990-02-02"},
		{"2017-02-03T10:40:30", "2017-02-03T10:40:30"},
		{"2017-02-03T10:40:30", "2019-02-03T10:40:30"},
	}

	for i, ex := range examples {
		// check x <= y is true
		expr := fmt.Sprintf("date[\"%s\"] <= date[\"%s\"]", ex[0], ex[1])
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != true {
			t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
		}

		// all examples meet x <= y, so if the Strings are not equal, then x < y, therefore y <= x must be false
		if ex[0] != ex[1] {
			expr := fmt.Sprintf("date[\"%s\"] <= date[\"%s\"]", ex[1], ex[0])
			v, err := expression.Eval(expr, nil)

			if err != nil {
				t.Fatalf("[%d] Could not evaluate: %v", i, err)
			}

			if v != false {
				t.Errorf("[%d] Expected 'false' but got '%v'", i, v)
			}
		}
	}
}

func TestDateGreaterOrEqual(t *testing.T) {
	examples := [][2]string{
		// lesser date, greater date
		{"1990-01-01", "1990-01-01"},
		{"1990-02-02", "1990-01-01"},
		{"2017-02-03T10:40:30", "2017-02-03T10:40:30"},
		{"2019-02-03T10:40:30", "2017-02-03T10:40:30"},
	}

	for i, ex := range examples {
		// check x <= y is true
		expr := fmt.Sprintf("date[\"%s\"] >= date[\"%s\"]", ex[0], ex[1])
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != true {
			t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
		}

		// all examples meet x >= y, so if the Strings are not equal, then x > y, therefore y >= x must be false
		if ex[0] != ex[1] {
			expr := fmt.Sprintf("date[\"%s\"] >= date[\"%s\"]", ex[1], ex[0])
			v, err := expression.Eval(expr, nil)

			if err != nil {
				t.Fatalf("[%d] Could not evaluate: %v", i, err)
			}

			if v != false {
				t.Errorf("[%d] Expected 'false' but got '%v'", i, v)
			}
		}
	}
}

func TestDateEqualAndNotEqual(t *testing.T) {
	// all examples should be different
	examples := []string{
		"1990-01-01", "1990-01-02", "2017-02-03T10:40:30", "2019-02-03T10:40:30",
	}

	for i, ex := range examples {
		for j := range examples {
			if i == j {
				// check x == x
				expr := fmt.Sprintf("date[\"%s\"] == date[\"%s\"]", ex, ex)
				v, err := expression.Eval(expr, nil)

				if err != nil {
					t.Fatalf("[%d] Could not evaluate: %v", i, err)
				}

				if v != true {
					t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
				}

				expr = fmt.Sprintf("date[\"%s\"] != date[\"%s\"]", ex, ex)
				v, err = expression.Eval(expr, nil)

				if err != nil {
					t.Fatalf("[%d] Could not evaluate: %v", i, err)
				}

				if v != false {
					t.Errorf("[%d] Expected 'false' but got '%v'", i, v)
				}
			} else {
				// check any two different examples are not equal
				expr := fmt.Sprintf("date[\"%s\"] != date[\"%s\"]", ex, examples[j])
				v, err := expression.Eval(expr, nil)

				if err != nil {
					t.Fatalf("[%d, %d] Could not evaluate: %v", i, j, err)
				}

				if v != true {
					t.Errorf("[%d, %d] Expected 'true' but got '%v': %s", i, j, v, expr)
				}

				expr = fmt.Sprintf("date[\"%s\"] == date[\"%s\"]", ex, examples[j])
				v, err = expression.Eval(expr, nil)

				if err != nil {
					t.Fatalf("[%d, %d] Could not evaluate: %v", i, j, err)
				}

				if v != false {
					t.Errorf("[%d, %d] Expected 'false' but got '%v': %s", i, j, v, expr)
				}
			}
		}
	}
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

	if result, ok := v.(*expression.DateTime); ok {
		expectedTime, err := time.Parse(format, expected)
		if err != nil {
			panic(err)
		}

		if result.Time != expectedTime {
			t.Errorf("Expected '%v' but got '%v'", expectedTime, v)
		}
	} else {
		t.Errorf("Expected DateTime, got '%v'", v)
	}
}
