package mg

import (
	"io"
	"log"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

type IfContent struct {
	Text      string
	condition *expression.Expression
	Location  *Location
	contents  []Content
	resolver  FileResolver
}

var _ Content = (*IfContent)(nil)
var _ ContentContainer = (*IfContent)(nil)

func NewIfInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	cond, err := expression.ParseExpr(arg)

	if err != nil {
		log.Printf("WARNING: (%s) Malformed if instruction: (%v)", location.String(), err)
		return unevaluatedExpression(original, location)
	}

	return &IfContent{
		Text:      original,
		condition: &cond,
		Location:  location,
		resolver:  resolver,
	}
}

func (ic *IfContent) AppendContent(content Content) {
	ic.contents = append(ic.contents, content)
}

func (ic *IfContent) GetLocation() *Location {
	return ic.Location
}

func (ic *IfContent) IsScoped() bool {
	return true
}

func (ic *IfContent) Write(writer io.Writer, context Context) ([]Content, error) {
	res, err := expression.EvalExpr(ic.condition, context)
	if err != nil {
		log.Printf("ERROR: If condition could not be evaluated: %v", err)
		return unevaluatedExpressions(ic.Text, ic.Location), nil
	}

	switch res {
	case true:
		return ic.contents, nil
	case false:
	case nil:
		return nil, nil
	default:
		log.Printf("INFO: If condition evaluated to non-boolean value, assuming false: %v", res)
	}

	return nil, nil
}
