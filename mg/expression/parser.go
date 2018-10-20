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

type Context interface {
	Get(name string) (interface{}, bool)
}

type MapContext struct {
	Map map[string]interface{}
}

func (m *MapContext) Get(name string) (interface{}, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m.Map[name]
	return v, ok
}

func ParseExpr(expr string) (Expression, error) {
	e, err := parser.ParseExpr(expr)
	if err != nil {
		return Expression{}, err
	}
	return Expression{expr: e}, nil
}

func Eval(expr string, context Context) (interface{}, error) {
	e, err := ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	return EvalExpr(e, context)
}

func EvalExpr(e Expression, contex Context) (interface{}, error) {
	return eval(e.expr, contex)
}

func eval(e ast.Expr, ctx Context) (interface{}, error) {
	switch ex := e.(type) {
	case *ast.BasicLit:
		return parseLiteral(ex.Value), nil
	case *ast.Ident:
		return resolveIdentifier(ex.Name, ctx)
	case *ast.BinaryExpr:
		return resolveBinaryExpr(ex.X, ex.Op, ex.Y, ctx)
	case *ast.CompositeLit:
		return resolveCompositeLit(ex, ctx)
	case *ast.ParenExpr:
		return eval(ex.X, ctx)
	case *ast.UnaryExpr:
		return resolveUnary(ex, ctx)
	case *ast.SelectorExpr:
		return resolveAccessField(ex, ctx)
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
	panic(fmt.Sprintf("Unrecognized literal: %s", s))
}

func resolveIdentifier(name string, ctx Context) (interface{}, error) {
	if name == "true" {
		return true, nil
	}
	if name == "false" {
		return false, nil
	}
	v, ok := ctx.Get(name)
	if ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("cannot resolve identifier: %s", name))
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
	case token.SUB:
		return subtract(xv, yv)
	case token.MUL:
		return multiply(xv, yv)
	case token.QUO:
		return divide(xv, yv)
	case token.REM:
		return rem(xv, yv)
	case token.EQL:
		return eq(xv, yv)
	case token.NEQ:
		return neq(xv, yv)
	case token.LSS:
		return lss(xv, yv)
	case token.GTR:
		return gtr(xv, yv)
	case token.LEQ:
		return leq(xv, yv)
	case token.GEQ:
		return geq(xv, yv)
	case token.LAND:
		return and(xv, yv)
	case token.LOR:
		return or(xv, yv)
	}
	return nil, errors.New(fmt.Sprintf("unknown operator %s", t))
}

func resolveCompositeLit(cl *ast.CompositeLit, ctx Context) (interface{}, error) {
	array := make([]interface{}, len(cl.Elts), len(cl.Elts))
	for i, v := range cl.Elts {
		item, err := eval(v, ctx)
		if err != nil {
			return nil, err
		}
		array[i] = item
	}
	return array, nil
}

func resolveUnary(expr *ast.UnaryExpr, ctx Context) (interface{}, error) {
	v, err := eval(expr.X, ctx)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case token.NOT:
		return not(v)
	case token.SUB:
		return minus(v)
	case token.ADD:
		return plus(v)
	}
	return nil, errors.New(fmt.Sprintf("operator %s cannot be used with %v",
		expr.Op, expr.X))
}

func resolveAccessField(expr *ast.SelectorExpr, ctx Context) (interface{}, error) {
	rcv, err := eval(expr.X, ctx)
	if err != nil {
		return nil, err
	}
	if v, ok := ToContext(rcv); ok {
		return eval(expr.Sel, v)
	}
	return nil, errors.New(fmt.Sprintf("cannot access properties of object: %v", rcv))
}

func ToContext(ctx interface{}) (Context, bool) {
	if v, ok := ctx.(Context); ok {
		return v, true
	}
	if v, ok := ctx.(map[string]interface{}); ok {
		return &MapContext{Map: v}, true
	}
	return nil, false
}
