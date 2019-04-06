package expression

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"time"
)

var defaultDateLayouts = []string{"2006-01-02T15:04:05", "2006-01-02T15:04", "2006-01-02"}

// Expression is a parsed Magnanimous expression.
//
// It can be evaluated with the [EvalExpr] function.
type Expression struct {
	expr ast.Expr
}

// DateTime is the result of evaluating a date expression.
type DateTime struct {
	Time   time.Time
	Format string
}

// Path is the result of evaluating a path expression.
type Path struct {
	Value string
}

type PathProperty struct {
	Path *Path
	Name string
}

// Context contains the bindings available for an expression.
type Context interface {
	Get(name string) (interface{}, bool)
}

// MapContext is a simple implementation of [Context].
type MapContext struct {
	Map map[string]interface{}
}

// Get the value of a binding from the context.
func (m *MapContext) Get(name string) (interface{}, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m.Map[name]
	return v, ok
}

// Get the value of a property of the file point at by the Path.
func (p *Path) Get(name string) (interface{}, bool) {
	return &PathProperty{Path: p, Name: name}, true
}

// ParseExpr parses the given string as a Magnanimous expression.
func ParseExpr(expr string) (Expression, error) {
	expr = strings.Trim(expr, " ")
	if strings.HasPrefix(expr, "[") && strings.HasSuffix(expr, "]") {
		expr = fmt.Sprintf("[]interface{}{%s}", expr[1:len(expr)-1])
	}
	e, err := parser.ParseExpr(expr)
	if err != nil {
		return Expression{}, err
	}
	return Expression{expr: e}, nil
}

// Eval evaluates the given string as a Magnanimous expression, making the context bindings
// available to the expression.
func Eval(expr string, context Context) (interface{}, error) {
	e, err := ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	return EvalExpr(e, context)
}

// EvalExpr evaluates the given Magnanimous expression, making the context bindings
//// available to the expression.
func EvalExpr(e Expression, context Context) (interface{}, error) {
	return eval(e.expr, context)
}

func eval(e ast.Expr, context Context) (interface{}, error) {
	switch ex := e.(type) {
	case *ast.BasicLit:
		return parseLiteral(ex.Value), nil
	case *ast.Ident:
		return resolveIdentifier(ex.Name, context)
	case *ast.BinaryExpr:
		return resolveBinaryExpr(ex.X, ex.Op, ex.Y, context)
	case *ast.CompositeLit:
		return resolveCompositeLit(ex, context)
	case *ast.ParenExpr:
		return eval(ex.X, context)
	case *ast.UnaryExpr:
		return resolveUnary(ex, context)
	case *ast.SelectorExpr:
		return resolveAccessField(ex, context)
	case *ast.IndexExpr:
		return resolveIndexExpr(ex, context)
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
	if ctx != nil {
		v, ok := ctx.Get(name)
		if ok {
			return v, nil
		}
	}

	switch name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return nil, nil
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
		return Equal(xv, yv)
	case token.NEQ:
		return NotEqual(xv, yv)
	case token.LSS:
		return Less(xv, yv)
	case token.GTR:
		return Greater(xv, yv)
	case token.LEQ:
		return LessOrEq(xv, yv)
	case token.GEQ:
		return GreaterOrEq(xv, yv)
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

func resolveIndexExpr(expr *ast.IndexExpr, ctx Context) (interface{}, error) {
	switch rcv := expr.X.(type) {
	case *ast.Ident:
		if rcv.Name == "date" {
			idx, err := eval(expr.Index, ctx)
			if err == nil {
				if d, ok := idx.(string); ok {
					return parseDate(d, "02 Jan 2006, 03:04 PM")
				} else {
					return nil, errors.New("malformed date expression (should be like date[\"2006-01-02T15:04:00\"])")
				}
			} else {
				return nil, err
			}
		}
		if rcv.Name == "path" {
			idx, err := eval(expr.Index, ctx)
			if err == nil {
				if p, ok := idx.(string); ok {
					return &Path{Value: p}, nil
				} else {
					return nil, errors.New("malformed path expression (should be like path[\"to/file.txt\"])")
				}
			} else {
				return nil, err
			}
		}
		return nil, fmt.Errorf("unsupported index expression (only date[] and path[] supported): %v", rcv)
	case *ast.IndexExpr:
		if i, ok := rcv.X.(*ast.Ident); ok {
			if i.Name == "date" {
				idx2, err := eval(expr.Index, ctx)
				if err == nil {
					if format, ok := idx2.(string); ok {
						idx1, err := eval(rcv.Index, ctx)
						if err == nil {
							if date, ok := idx1.(string); ok {
								return parseDate(date, format)
							}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				return nil, errors.New("malformed date expression (should be like date[\"2006-01-02T15:04:00\"][layout])")
			}
		}
		return nil, errors.New("unsupported index expression (only date[][] supported)")
	}

	return nil, errors.New("malformed index expression (only date[][] supported)")
}

func parseDate(idx string, format string) (*DateTime, error) {
	if idx == "now" { // special case
		return &DateTime{Format: format, Time: time.Now()}, nil
	}
	for _, layout := range defaultDateLayouts {
		date, err := time.Parse(layout, idx)
		if err == nil {
			return &DateTime{Format: format, Time: date}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("invalid date: %v (valid formats: %v)", idx, defaultDateLayouts))
}

// ToContext attempts to convert a variable to a [Context].
func ToContext(ctx interface{}) (Context, bool) {
	if v, ok := ctx.(Context); ok {
		return v, true
	}
	if v, ok := ctx.(map[string]interface{}); ok {
		return &MapContext{Map: v}, true
	}
	return nil, false
}
