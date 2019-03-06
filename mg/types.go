package mg

import (
	"bytes"
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

// ContextStack is a stack of InclusionChainItems.
//
// Used to keep state when writing nested Content.
type ContextStack struct {
	chain []ContextStackItem
}

// ContextStackItem is an object that contains the local Context and Location.
//
// It is used by Content implementations to write their contents using the given context to resolve data.
// Content implementations that create new scopes must add an item to the ContextStack.
type ContextStackItem struct {
	Location *Location
	Context  Context
}

// FileResolver defines how Magnanimous finds source files.
type FileResolver interface {
	// FilesIn return the files in a certain directory, or an error if something goes wrong.
	FilesIn(dir string, from *Location) (dirPath string, f []WebFile, e error)
	// Resolve resolves a path given a location to resolve it from.
	// It allows Magnanimous to resolve relative paths correctly.
	Resolve(path string, from *Location) string
}

// ContentContainer is a collection of Content.
//
// Content implementations that have nested Content must implement this interface.
type ContentContainer interface {
	GetContents() []Content
	AppendContent(content Content)
}

// Content is a processed unit of a source file.
//
// Implementations of Content define how Magnanimous instructions behave.
// A Content may contain nested Content parts, in which case it implements ContentContainer.
type Content interface {
	// Write writes its contents using the given writer.
	//
	// The files argument contains all parsed source files.
	//
	// The stack contains context in which local data can be stored.
	// Each implementation of Content that starts a new scope must push a new item onto the stack.
	Write(writer io.Writer, files WebFilesMap, stack ContextStack) error
}

// Context represents the current context of a Content being written.
type Context interface {
	Get(name string) (interface{}, bool)
	Set(name string, value interface{})
	IsEmpty() bool
}

// ProcessedFile is the result of parsing a source file.
type ProcessedFile struct {
	contents     []Content
	NewExtension string
	Path         string
}

var _ ContentContainer = (*ProcessedFile)(nil)

func (f *ProcessedFile) GetContents() []Content {
	return f.contents
}

func (f *ProcessedFile) AppendContent(content Content) {
	f.contents = append(f.contents, content)
}

// ResolveContext evaluates all of the [DefineContent] instructions at the top-level scope
// of the [ProcessedFile].
func (f *ProcessedFile) ResolveContext(files WebFilesMap, stack ContextStack) Context {
	ctx := make(map[string]interface{})
	for _, c := range f.contents {
		if content, ok := c.(*DefineContent); ok {
			v, ok := content.Eval(files, stack)
			if ok {
				ctx[content.Name] = v
			}
		}
	}
	return &mapContext{ctx}
}

// Bytes returns the bytes of the processed file.
func (f *ProcessedFile) Bytes(files WebFilesMap, stack ContextStack) ([]byte, error) {
	return body(f, files, stack)
}

func body(c ContentContainer, files WebFilesMap, stack ContextStack) ([]byte, error) {
	var b bytes.Buffer
	b.Grow(512)
	for _, c := range c.GetContents() {
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
