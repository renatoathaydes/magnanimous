package mg

import (
	"io"
	"log"
	"strings"
)

type SlotContent struct {
	UnscopedContent
	Name     string
	Text     string
	Location *Location
	contents []Content
}

var _ Content = (*SlotContent)(nil)
var _ ContentContainer = (*SlotContent)(nil)
var _ Definition = (*SlotContent)(nil)

// slotEval is a Content that gets written when a slot is evaluated.
type slotEval struct {
	location *Location
	contents []Content
}

var _ Content = (*slotEval)(nil)

func NewSlotInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 1 {
		variable := parts[0]
		return &SlotContent{Name: variable, Text: original, Location: location}
	}
	log.Printf("WARNING: (%s) malformed slot instruction: %s", location.String(), arg)
	return unevaluatedExpression(original, location)
}

func (s *SlotContent) GetName() string {
	return s.Name
}

func (s *SlotContent) AppendContent(content Content) {
	s.contents = append(s.contents, content)
}

func (s *SlotContent) GetLocation() *Location {
	return s.Location
}

func (s *SlotContent) Write(writer io.Writer, context Context) ([]Content, error) {
	if v, ok := s.Eval(context); ok {
		context.Set(s.Name, v)
		return nil, nil
	}
	return unevaluatedExpressions(s.Text, s.Location), nil
}

func (s *SlotContent) Eval(context Context) (interface{}, bool) {
	return &slotEval{location: s.Location, contents: s.contents}, true
}

func (s *slotEval) GetLocation() *Location {
	return s.location
}

func (s *slotEval) IsScoped() bool {
	return true
}

func (s *slotEval) Write(writer io.Writer, context Context) ([]Content, error) {
	return s.contents, nil
}
