package mg

import (
	"fmt"
	"io"
	"log"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

// EvalContent is a MultiContent that represents a Magnanimous Language Expression.
type EvalContent struct {
	UnscopedContent
	Expr     *expression.Expression
	Text     string
	Location *Location
	resolver FileResolver
}

var _ Content = (*EvalContent)(nil)

func NewEvalInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	expr, err := expression.ParseExpr(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original, location)
	}
	return &EvalContent{Expr: &expr, Location: location, Text: original, resolver: resolver}
}

func (e *EvalContent) GetLocation() *Location {
	return e.Location
}

func (e *EvalContent) Write(writer io.Writer, context Context) ([]Content, error) {
	v, err := e.eval(context)
	if err == nil {
		switch c := v.(type) {
		case Content:
			return []Content{c}, nil
		case string:
			if _, err := writer.Write([]byte(c)); err != nil {
				return nil, err
			}
			return nil, nil
		case nil:
			return nil, nil
		}
		err = fmt.Errorf("value has unexpected type (should be string or Content): %v", v)
	}
	log.Printf("ERROR: (%s) eval failure [%s]: %s", e.Location.String(), e.Text, err.Error())
	return unevaluatedExpressions(e.Text, e.Location), nil
}

// eval evaluates the given expression to either a string or a Content, or an error.
func (e *EvalContent) eval(context Context) (interface{}, error) {
	r, err := expression.EvalExpr(e.Expr, context)
	if err == nil {
		// an expression can evaluate to Content, such as a slot
		if c, ok := r.(Content); ok {
			return c, nil
		}
		// evaluate special types to a simple string to write
		return e.evalSpecialType(r, context)
	}
	return nil, err
}

func (e *EvalContent) evalSpecialType(r interface{}, context Context) (string, error) {
	switch v := r.(type) {
	case nil:
		return "", nil
	case *expression.DateTime:
		if v.Time != nil {
			return v.Time.Format(v.Format), nil
		}
		// DateTime does not have a Time, it has a Path instead
		f, err := getInclusionByPath(e.toInclusion(v.Path), e.resolver, context, false)
		if err != nil {
			return "", err
		}
		return f.Processed.LastUpdated.Format(v.Format), nil
	case *expression.Path:
		f, err := getInclusionByPath(e.toInclusion(v), e.resolver, context, false)
		if err != nil {
			return "", err
		}
		return f.Processed.Path, nil
	case *expression.PathProperty:
		f, err := getInclusionByPath(e.toInclusion(v.Path), e.resolver, context, false)
		if err != nil {
			return "", err
		}
		ctx := f.Processed.ResolveContext(context, false)
		if prop, ok := ctx.Get(v.Name); ok {
			return e.evalSpecialType(prop, context)
		}
		return "", fmt.Errorf("File at path %s has no such property: %s", v.Path.Value, v.Name)
	}
	// no special type found, stringify it
	return fmt.Sprintf("%v", r), nil
}

func (e *EvalContent) toInclusion(p *expression.Path) Inclusion {
	inc := pathInclusion{e.Location, p}
	return &inc
}

func (e *EvalContent) String() string {
	return fmt.Sprintf("EvalContent{%s}", e.Text)
}

type pathInclusion struct {
	location *Location
	path     *expression.Path
}

var _ Inclusion = (*pathInclusion)(nil)

func (p *pathInclusion) GetLocation() *Location {
	return p.location
}

func (p *pathInclusion) GetPath() string {
	return p.path.Value
}
