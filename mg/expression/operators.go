package expression

import (
	"errors"
	"fmt"
)

type numberOp func(x float64, y float64) interface{}
type stringOp func(x string, y string) interface{}
type boolOp func(x bool, y bool) interface{}
type opOther func() (interface{}, error)

func op(x interface{}, y interface{}, no numberOp, so stringOp, oe opOther) (interface{}, error) {
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
	return oe()
}

func bop(x interface{}, y interface{}, bo boolOp, oe opOther) (interface{}, error) {
	if bo != nil {
		xf, ok := x.(bool)
		if ok {
			yf, ok := y.(bool)
			if ok {
				return bo(xf, yf), nil
			}
		}
	}
	return oe()
}

func add(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x + y
	}, func(x string, y string) interface{} {
		return x + y
	}, func() (interface{}, error) {
		xs := fmt.Sprintf("%v", x)
		ys := fmt.Sprintf("%v", y)
		return xs + ys, nil
	})
}

func subtract(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x - y
	}, nil, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot subtract %v from %v", y, x))
	})
}

func multiply(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x * y
	}, nil, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot multiply %v and %v", x, y))
	})
}

func divide(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x / y
	}, nil, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot divide %v by %v", x, y))
	})
}

func rem(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return float64(int64(x) % int64(y))
	}, nil, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot divide %v by %v (remainder)", x, y))
	})
}

func Equal(x interface{}, y interface{}) (interface{}, error) {
	return x == y, nil
}

func NotEqual(x interface{}, y interface{}) (interface{}, error) {
	return x != y, nil
}

func Less(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x < y
	}, func(x string, y string) interface{} {
		return x < y
	}, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func Greater(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x > y
	}, func(x string, y string) interface{} {
		return x > y
	}, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func LessOrEq(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x <= y
	}, func(x string, y string) interface{} {
		return x <= y
	}, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func GreaterOrEq(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x >= y
	}, func(x string, y string) interface{} {
		return x >= y
	}, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot compare %v and %v", x, y))
	})
}

func and(x interface{}, y interface{}) (interface{}, error) {
	return bop(x, y, func(x bool, y bool) interface{} {
		return x && y
	}, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot AND %v and %v", x, y))
	})
}

func or(x interface{}, y interface{}) (interface{}, error) {
	if x == nil || x == float64(0) || x == false || x == "" {
		return y, nil
	}
	return x, nil
}

func not(x interface{}) (interface{}, error) {
	return bop(x, true, func(x bool, y bool) interface{} {
		return !x
	}, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot negate %v", x))
	})
}

func minus(x interface{}) (interface{}, error) {
	return op(x, float64(0), func(x float64, y float64) interface{} {
		return -x
	}, nil, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot use '-' on %v", x))
	})
}

func plus(x interface{}) (interface{}, error) {
	return op(x, float64(0), func(x float64, y float64) interface{} {
		return x
	}, nil, func() (interface{}, error) {
		return nil, errors.New(fmt.Sprintf("cannot use '+' on %v", x))
	})
}
