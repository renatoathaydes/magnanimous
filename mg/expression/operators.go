package expression

import (
	"errors"
	"fmt"
)

type numberOp func(x float64, y float64) interface{}
type stringOp func(x string, y string) interface{}
type opErr func() error

func op(x interface{}, y interface{}, no numberOp, so stringOp, oe opErr) (interface{}, error) {
	if no != nil {
		xf, ok := x.(float64)
		if ok {
			yf, ok := y.(float64)
			if ok {
				return no(xf, yf), nil
			}
		}
	}
	if so != nil {
		xs, ok := x.(string)
		if ok {
			ys, ok := y.(string)
			if ok {
				return so(xs, ys), nil
			}
		}
	}
	return nil, oe()
}

func add(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x + y
	}, func(x string, y string) interface{} {
		return x + y
	}, func() error {
		return errors.New(fmt.Sprintf("cannot add %v to %v", x, y))
	})
}

func subtract(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x - y
	}, nil, func() error {
		return errors.New(fmt.Sprintf("cannot subtract %v from %v", y, x))
	})
}

func multiply(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x * y
	}, nil, func() error {
		return errors.New(fmt.Sprintf("cannot multiply %v and %v", x, y))
	})
}

func divide(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x / y
	}, nil, func() error {
		return errors.New(fmt.Sprintf("cannot divide %v by %v", x, y))
	})
}

func eq(x interface{}, y interface{}) (interface{}, error) {
	return x == y, nil
}

func neq(x interface{}, y interface{}) (interface{}, error) {
	return x != y, nil
}

func lss(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x < y
	}, func(x string, y string) interface{} {
		return x < y
	}, func() error {
		return errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func gtr(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x > y
	}, func(x string, y string) interface{} {
		return x > y
	}, func() error {
		return errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func leq(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x <= y
	}, func(x string, y string) interface{} {
		return x <= y
	}, func() error {
		return errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func geq(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x >= y
	}, func(x string, y string) interface{} {
		return x >= y
	}, func() error {
		return errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}
