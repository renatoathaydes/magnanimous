package mg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Magnanimous is the entry point of the magnanimous library.
//
// It can be used to read a source directory via the ReadAll() function.
type Magnanimous struct {
	// SourcesDir is the directory containing Magnanimous' source code.
	SourcesDir string
}

// WebFilesMap contains the result of reading a source directory with ReadAll().
type WebFilesMap struct {
	GlobalContext RootScope
	// WebFiles is a map from each file path to the parsed WebFile.
	WebFiles map[string]WebFile
}

// WebFile is a parsed source file.
type WebFile struct {
	BasePath    string
	Name        string
	Processed   *ProcessedFile
	NonWritable bool
}

// Location is used to show where in the source code errors or warning messages originate.
type Location struct {
	Origin string
	Row    uint32
	Col    uint32
}

// InclusionChainItem is an object that contains the current scope.
//
// It is used by Content implementations to write their contents using the given scope to resolve data.
type InclusionChainItem struct {
	Location *Location
	scope    Scope
}

// FileResolver defines how Magnanimous finds source files.
type FileResolver interface {
	// FilesIn return the files in a certain directory, or an error if something goes wrong.
	FilesIn(dir string, from Location) (dirPath string, f []WebFile, e error)
	// Resolve resolves a path given a location to resolve it from.
	// It allows Magnanimous to resolve relative paths correctly.
	Resolve(path string, from Location) string
}

// ContentContainer is a collection of Content.
//
// Content implementations that have nested Content must implement this interface.
type ContentContainer interface {
	GetContents() []Content
}

// Content is a processed unit of a source file.
//
// Implementations of Content define how Magnanimous instructions behave.
// A Content may contain nested Content parts, in which case it implements ContentContainer.
type Content interface {
	// Write writes its contents using the given writer.
	//
	// The files argument contains all parsed source files.
	// The inclusionChain contains the scope under which the Content is being written.
	//
	// When writing nested contents, implementations must append their immediate InclusionChainItem
	// to the slice before calling the Write method on the nested Content instances.
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
