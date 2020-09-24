package mg

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"
)

// ProcessedFile is the result of parsing a source file.
type ProcessedFile struct {
	contents     []Content
	NewExtension string
	BasePath     string
	Path         string
	LastUpdated  time.Time
}

var _ ContentContainer = (*ProcessedFile)(nil)

func (f *ProcessedFile) AppendContent(content Content) {
	f.contents = append(f.contents, content)
}

func (f *ProcessedFile) GetLocation() *Location {
	loc := Location{Origin: f.Path, Row: 0, Col: 0}
	return &loc
}

func (f *ProcessedFile) GetContents() []Content {
	return f.contents
}

// ResolveContext evaluates all of the [Definition]s at the top-level scope
// of the [ProcessedFile].
//
// If [inPlace] is true, the definitions are set into the given context, which is also returned.
// Otherwise, a new context is created and all variables are set into it, then this new context is returned.
func (f *ProcessedFile) ResolveContext(context Context, inPlace bool) Context {
	var ctx Context
	if inPlace {
		ctx = context
	} else {
		ctx = newFileContext(f, context)
	}
	resolveContext(f.GetContents(), ctx)
	return ctx
}

func resolveContext(contents []Content, context Context) {
	var buffer bytes.Buffer
	for _, c := range contents {
		if def, ok := c.(Definition); ok {
			_, err := def.Write(&buffer, context)
			if err != nil {
				log.Printf("ERROR: (%s) eval failure [%s]: %s", def.GetLocation().String(), def.GetName(), err.Error())
			}
		}
	}
}

func (f *ProcessedFile) String() string {
	var contentsBuilder strings.Builder
	contentsBuilder.WriteString("[ ")
	for _, c := range f.contents {
		contentsBuilder.WriteString(fmt.Sprintf("%T ", c))
	}
	contentsBuilder.WriteString("]")
	return fmt.Sprintf("ProcessedFile{%s, %s}", contentsBuilder.String(), f.NewExtension)
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Origin, l.Row, l.Col)
}
