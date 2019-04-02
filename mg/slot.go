package mg

import (
	"io"
	"log"
	"strings"
)

type SlotContent struct {
	Name     string
	Text     string
	Location *Location
	contents []Content
}

var _ Content = (*SlotContent)(nil)
var _ ContentContainer = (*SlotContent)(nil)

func NewSlotInstruction(arg string, location *Location, original string) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 1 {
		variable := parts[0]
		return &SlotContent{Name: variable, Location: location}
	}
	log.Printf("WARNING: (%s) malformed slot instruction: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func (s *SlotContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	// slot does not write anything, it must be evaluated to resolve its body
	stack.Top().Set(s.Name, s)
	return nil
}

func (s *SlotContent) GetContents() []Content {
	if isMd(s.Location.Origin) {
		return []Content{&HtmlFromMarkdownContent{MarkDownContent: s.contents}}
	} else {
		return s.contents
	}
}

func (s *SlotContent) AppendContent(content Content) {
	s.contents = append(s.contents, content)
}
