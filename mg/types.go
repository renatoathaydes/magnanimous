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

type WebFilesMap struct {
	GlobalContext RootScope
	WebFiles      map[string]WebFile
}

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

type SideEffectContent interface {
	Run(files *WebFilesMap, inclusionChain []InclusionChainItem)
}

type StringContent struct {
	Text string
}

type RootScope Context

type Context interface {
	Get(name string) (interface{}, bool)
	Set(name string, value interface{})
	IsEmpty() bool
	mixInto(other Context)
}

type Scope interface {
	AppendContent(content Content)
	Context() Context
	Parent() Scope
}

type MapContext struct {
	Map map[string]interface{}
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
var _ Context = (*MapContext)(nil)

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

func (f *ProcessedFile) Context() Context {
	if f.context == nil {
		f.context = make(map[string]interface{})
	}
	return &MapContext{Map: f.context}
}

func (f *ProcessedFile) Parent() Scope {
	return nil
}

func (f *ProcessedFile) AppendContent(content Content) {
	s := len(f.scopeStack)
	var topScope Scope = f
	if s > 0 {
		topScope = f.scopeStack[s-1]
		topScope.AppendContent(content)
	} else {
		f.contents = append(f.contents, content)
	}
	if newScope, ok := content.(Scope); ok {
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

func (m *MapContext) Get(name string) (interface{}, bool) {
	v, ok := m.Map[name]
	return v, ok
}

func (m *MapContext) Set(name string, value interface{}) {
	m.Map[name] = value
}

func (m *MapContext) IsEmpty() bool {
	return len(m.Map) == 0
}

func (m *MapContext) mixInto(other Context) {
	for k, v := range m.Map {
		other.Set(k, v)
	}
}
