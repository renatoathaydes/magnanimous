package expression

import (
	"fmt"
	"reflect"
)

type numberOp func(x float64, y float64) interface{}
type stringOp func(x string, y string) interface{}
type boolOp func(x bool, y bool) interface{}
type timeOp func(x *DateTime, y *DateTime) interface{}
type opOther func() (interface{}, error)

func op(x interface{}, y interface{}, no numberOp, so stringOp, to timeOp, oe opOther) (interface{}, error) {
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
	if to != nil {
		xt, ok := x.(*DateTime)
		if ok {
			yt, ok := y.(*DateTime)
			if ok {
				return to(xt, yt), nil
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
	}, nil, func() (interface{}, error) {
		xs := fmt.Sprintf("%v", x)
		ys := fmt.Sprintf("%v", y)
		return xs + ys, nil
	})
}

func subtract(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x - y
	}, nil, nil, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot subtract %v from %v", y, x)
	})
}

func multiply(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x * y
	}, nil, nil, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot multiply %v and %v", x, y)
	})
}

func divide(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x / y
	}, nil, nil, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot divide %v by %v", x, y)
	})
}

func rem(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return float64(int64(x) % int64(y))
	}, nil, nil, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot divide %v by %v (remainder)", x, y)
	})
}

func Equal(x interface{}, y interface{}) (interface{}, error) {
	return reflect.DeepEqual(x, y), nil
}

func NotEqual(x interface{}, y interface{}) (interface{}, error) {
	return !reflect.DeepEqual(x, y), nil
}

func Less(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x < y
	}, func(x string, y string) interface{} {
		return x < y
	}, func(x *DateTime, y *DateTime) interface{} {
		if x.Time != nil && y.Time != nil {
			return x.Time.Before(*y.Time)
		}
		return fmt.Errorf("cannot compare file time instances")
	}, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot compare %v and %v", x, y)
	})
}

func Greater(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x > y
	}, func(x string, y string) interface{} {
		return x > y
	}, func(x *DateTime, y *DateTime) interface{} {
		if x.Time != nil && y.Time != nil {
			return x.Time.After(*y.Time)
		}
		return fmt.Errorf("cannot compare file time instances")
	}, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot compare %v and %v", x, y)
	})
}

func LessOrEq(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x <= y
	}, func(x string, y string) interface{} {
		return x <= y
	}, func(x *DateTime, y *DateTime) interface{} {
		if x.Time != nil && y.Time != nil {
			return x.Time.Before(*y.Time) || *x.Time == *y.Time
		}
		return fmt.Errorf("cannot compare file time instances")
	}, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot compare %v and %v", x, y)
	})
}

func GreaterOrEq(x interface{}, y interface{}) (interface{}, error) {
	return op(x, y, func(x float64, y float64) interface{} {
		return x >= y
	}, func(x string, y string) interface{} {
		return x >= y
	}, func(x *DateTime, y *DateTime) interface{} {
		if x.Time != nil && y.Time != nil {
			return x.Time.After(*y.Time) || *x.Time == *y.Time
		}
		return fmt.Errorf("cannot compare file time instances")
	}, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot compare %v and %v", x, y)
	})
}

func and(x interface{}, y interface{}) (interface{}, error) {
	return bop(x, y, func(x bool, y bool) interface{} {
		return x && y
	}, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot AND %v and %v", x, y)
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
		return nil, fmt.Errorf("cannot negate %v", x)
	})
}

func minus(x interface{}) (interface{}, error) {
	return op(x, float64(0), func(x float64, y float64) interface{} {
		return -x
	}, nil, nil, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot use '-' on %v", x)
	})
}

func plus(x interface{}) (interface{}, error) {
	return op(x, float64(0), func(x float64, y float64) interface{} {
		return x
	}, nil, nil, func() (interface{}, error) {
		return nil, fmt.Errorf("cannot use '+' on %v", x)
	})
}
