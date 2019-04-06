package mg

import (
	"bytes"
	"fmt"
	"strings"
)

func (f *ProcessedFile) GetContents() []Content {
	return f.contents
}

func (f *ProcessedFile) AppendContent(content Content) {
	f.contents = append(f.contents, content)
}

// ResolveContext evaluates all of the [DefineContent] instructions at the top-level scope
// of the [ProcessedFile].
func (f *ProcessedFile) ResolveContext(files WebFilesMap, stack ContextStack) Context {
	ctx := newFileContext(f)
	stack = stack.PushContext(ctx)
	for _, c := range f.expandedContents() {
		if content, ok := c.(*DefineContent); ok {
			v, ok := content.Eval(files, stack)
			if ok {
				ctx.Set(content.Name, v)
			}
		}
	}
	return ctx
}

// Bytes returns the bytes of the processed file.
func (f *ProcessedFile) Bytes(files WebFilesMap, stack ContextStack) ([]byte, error) {
	return body(f, files, stack)
}

func (f *ProcessedFile) expandedContents() []Content {
	// expand contents in case this is markdown
	if c, ok := unwrapMarkdownContent(f); ok {
		return c
	}
	return f.contents
}

func body(c ContentContainer, files WebFilesMap, stack ContextStack) ([]byte, error) {
	return asBytes(c.GetContents(), files, stack)
}

func asBytes(c []Content, files WebFilesMap, stack ContextStack) ([]byte, error) {
	var b bytes.Buffer
	b.Grow(512)
	for _, c := range c {
		if c != nil {
			err := c.Write(&b, files, stack)
			if err != nil {
				return nil, err
			}
		}
	}
	return b.Bytes(), nil
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
