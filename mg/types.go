package mg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Magnanimous struct {
	SourcesDir string
}

type WebFilesMap map[string]WebFile

type WebFile struct {
	BasePath    string
	Name        string
	Processed   *ProcessedFile
	NonWritable bool
}

type Location struct {
	Origin string
	Row    uint32
	Col    uint32
}

type InclusionChainItem struct {
	Location *Location
	scope    Scope
}

type FileResolver interface {
	FilesIn(dir string, from Location) (dirPath string, f []WebFile, e error)
	Resolve(path string, from Location) string
}

type ContentContainer interface {
	GetContents() []Content
}

type Content interface {
	Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error
}

type StringContent struct {
	Text string
}

type RootScope map[string]interface{}

type Scope interface {
	AppendContent(content Content)
	Context() map[string]interface{}
	Parent() Scope
	setParent(scope Scope)
}

type ProcessedFile struct {
	contents     []Content
	scopeStack   []Scope
	context      map[string]interface{}
	rootScope    RootScope
	NewExtension string
}

var _ Scope = (*ProcessedFile)(nil)
var _ ContentContainer = (*ProcessedFile)(nil)

func (f *ProcessedFile) GetContents() []Content {
	return f.contents
}

// currentScope returns the current scope during parsing.
func (f *ProcessedFile) currentScope() Scope {
	s := len(f.scopeStack)
	if s > 0 {
		return f.scopeStack[s-1]
	}
	return f
}

func (f *ProcessedFile) Context() map[string]interface{} {
	if f.context == nil {
		f.context = make(map[string]interface{})
	}
	return f.context
}

func (f *ProcessedFile) Parent() Scope {
	return f.rootScope
}

func (f *ProcessedFile) setParent(content Scope) {
	switch root := content.(type) {
	case RootScope:
		f.rootScope = root
	default:
		panic("Cannot set parent on ProcessedFile as it's the top content scope")
	}
}

func (f *ProcessedFile) AppendContent(content Content) {
	s := len(f.scopeStack)
	var topScope Scope
	if s > 0 {
		topScope = f.scopeStack[s-1]
		topScope.AppendContent(content)
	} else {
		f.contents = append(f.contents, content)
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

func (f *ProcessedFile) Bytes(files WebFilesMap, inclusionChain []InclusionChainItem) ([]byte, error) {
	var b bytes.Buffer
	b.Grow(512)
	for _, c := range f.contents {
		if c != nil {
			err := c.Write(&b, files, inclusionChain)
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
	return fmt.Sprintf("ProcessedFile{%s, %v, %s}", contentsBuilder.String(), f.context, f.NewExtension)
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Origin, l.Row, l.Col)
}

var _ Scope = (*RootScope)(nil)

func (RootScope) AppendContent(content Content) {
	panic("RootScope cannot append content")
}

func (r RootScope) Context() map[string]interface{} {
	return r
}

func (RootScope) Parent() Scope {
	return nil
}

func (RootScope) setParent(scope Scope) {
	panic("Cannot set RootScope's parent")
}
