package mg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type WebFilesMap map[string]WebFile

type WebFile struct {
	BasePath    string
	Processed   *ProcessedFile
	NonWritable bool
}

type Location struct {
	Origin string
	Row    uint32
	Col    uint32
}

type Content interface {
	Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError
	IsMarkDown() bool
}

type StringContent struct {
	Text     string
	MarkDown bool
}

type HtmlFromMarkdownContent struct {
	MarkDownContent Content
}

type Scope interface {
	AppendContent(content Content)
	Context() map[string]interface{}
	Parent() Scope
	setParent(scope Scope)
}

type ScopeSensitive interface {
	setScope(scope Scope)
}

type ProcessedFile struct {
	Contents     []Content
	scopeStack   []Scope
	context      map[string]interface{}
	NewExtension string
}

var _ Scope = (*ProcessedFile)(nil)

// currentScope returns the current scope during parsing.
func (f *ProcessedFile) currentScope() Scope {
	s := len(f.scopeStack)
	if s > 0 {
		return f.scopeStack[s-1]
	}
	return f
}

func (f *ProcessedFile) Context() map[string]interface{} {
	return f.context
}

func (f *ProcessedFile) Parent() Scope {
	return nil
}

func (f *ProcessedFile) setParent(content Scope) {
	panic("Cannot set parent on ProcessedFile as it's the root content scope")
}

func (f *ProcessedFile) AppendContent(content Content) {
	s := len(f.scopeStack)
	var topScope Scope
	if s > 0 {
		topScope = f.scopeStack[s-1]
		topScope.AppendContent(content)
		h, ok := content.(ScopeSensitive)
		if ok {
			h.setScope(topScope)
		}
	} else {
		f.Contents = append(f.Contents, content)
		h, ok := content.(ScopeSensitive)
		if ok {
			h.setScope(f)
		}
	}
	newScope, ok := content.(Scope)
	if ok {
		newScope.setParent(topScope)
		f.scopeStack = append(f.scopeStack, newScope)
	}
}

func (f *ProcessedFile) EndScope() error {
	s := len(f.scopeStack)
	if s > 0 {
		f.scopeStack = f.scopeStack[0 : s-1]
		return nil
	} else {
		return errors.New("'end' does not match any previous instruction")
	}
}

func (f *ProcessedFile) Bytes(files WebFilesMap, inclusionChain []Location) []byte {
	var b bytes.Buffer
	b.Grow(512)
	for _, c := range f.Contents {
		if c != nil {
			c.Write(&b, files, inclusionChain)
		}
	}
	return b.Bytes()
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Origin, l.Row, l.Col)
}
