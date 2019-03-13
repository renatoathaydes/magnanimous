package mg

import (
	"io"
)

// Magnanimous is the entry point of the magnanimous library.
//
// It can be used to read a source directory via the ReadAll() function.
type Magnanimous struct {
	// SourcesDir is the directory containing Magnanimous' source code.
	SourcesDir string
	// Location of the global context relative to the "processed" directory.
	GlobalContex string
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
	locations []Location
	contexts  []Context
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
	// Write contents using the given writer.
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
	Remove(name string) interface{}
	IsEmpty() bool
}

// ProcessedFile is the result of parsing a source file.
type ProcessedFile struct {
	contents     []Content
	NewExtension string
	Path         string
}

var _ ContentContainer = (*ProcessedFile)(nil)
