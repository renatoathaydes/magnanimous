package expression

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type Expression struct {
	expr ast.Expr
}

type Context map[string]interface{}

func ParseExpr(expr string) (Expression, error) {
	e, err := parser.ParseExpr(expr)
	if err != nil {
		return Expression{}, err
	}
	return Expression{expr: e}, nil
}

func EvalExpr(e Expression, contex Context) (interface{}, error) {
	return eval(e.expr, contex)
}

func eval(e ast.Expr, ctx Context) (interface{}, error) {
	switch ex := e.(type) {
	case *ast.BasicLit:
		return parseLiteral(ex.Value), nil
	case *ast.BinaryExpr:
		return resolveBinaryExpr(ex.X, ex.Op, ex.Y, ctx)
	}
	return nil, errors.New(fmt.Sprintf("Unrecognized expression: %s", e))
}
func parseLiteral(s string) interface{} {
	if (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`")) {
		return s[1 : len(s)-1]
	}
	f, err := strconv.ParseFloat(s, 32)
	if err == nil {
		return f
	}
	b, err := strconv.ParseBool(s)
	if err == nil {
		return b
	}
	panic(fmt.Sprintf("Uncovered literal: %s", s))
}

func resolveBinaryExpr(x ast.Expr, t token.Token, y ast.Expr, ctx Context) (interface{}, error) {
	xv, err := eval(x, ctx)
	if err != nil {
		return nil, err
	}
	yv, err := eval(y, ctx)
	if err != nil {
		return nil, err
	}
	switch t {
	case token.ADD:
		return add(xv, yv)
	}
	return nil, errors.New("unknown operator")
}

func add(x interface{}, y interface{}) (interface{}, error) {
	xf, ok := x.(float64)
	if ok {
		yf, ok := y.(float64)
		if ok {
			return xf + yf, nil
		}
	}
	xs, ok := x.(string)
	if ok {
		ys, ok := y.(string)
		if ok {
			return xs + ys, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("cannot add %v to %v", x, y))
}
