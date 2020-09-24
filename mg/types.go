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
	BasePath       string
	Name           string
	Processed      *ProcessedFile
	NonWritable    bool
	SkipIfUpToDate bool
}

// Location is used to show where in the source code errors or warning messages originate.
type Location struct {
	Origin string
	Row    uint32
	Col    uint32
}

// FileResolver defines how Magnanimous finds source files.
type FileResolver interface {
	// FilesIn return the files in a certain directory, or an error if something goes wrong.
	FilesIn(dir string, from *Location) (dirPath string, f []WebFile, e error)
	// Resolve a path given a location to resolve it from.
	// It allows Magnanimous to resolve relative paths correctly.
	// To support up-paths, the [at] location calling the resolution (which can be different from [from])
	// is also needed.
	Resolve(path string, from *Location, at *Location) string
	// Get the file at the given path.
	//
	// This method does not attempt to resolve the path. If that is required, Resolve() should be called first.
	Get(path string) (*WebFile, bool)
}

// ContentContainer is a mutable collection of Content.
//
// Content implementations that have nested Content should implement this interface if they receive their
// nested contents during parsing.
type ContentContainer interface {

	// AppendContent adds nested content to this container.
	AppendContent(content Content)
}

// Definition is a Content that defines a new value within its context.
type Definition interface {
	Content

	// GetName returns the name of this definition.
	GetName() string

	// Eval evaluates this definition, inserting its name into the given context.
	//
	// Returns the evaluated definition's value and true if the evaluation succeeds, or
	// an undefined object and false otherwise.
	//
	// Eval does not insert the definition into the given Context, the caller is expected to do that.
	Eval(context Context) (interface{}, bool)
}

// Content is a processed unit of a source file.
//
// Implementations of Content define how Magnanimous instructions behave.
//
// If a Content returns other Contents when its [Write] method is called, the returned contents
// are written immediately, in a new scope if [IsScoped] returns true.
type Content interface {
	// GetLocation() returns content's original location.
	GetLocation() *Location

	// IsScoped returns true if this Content starts a new scope, false otherwise.
	IsScoped() bool

	// Write contents using the given writer.
	//
	// The stack contains context in which local data can be stored.
	//
	// Write should return Contents in cases where its evaluation results in yet more Contents
	// (e.g. the `include` instruction should return the included contents).
	Write(writer io.Writer, context Context) ([]Content, error)
}

// UnscopedContent is a base struct for unscoped Contents.
type UnscopedContent struct{}

// IsScoped returns false (see [Content]).
func (UnscopedContent) IsScoped() bool {
	return false
}

// Inclusion is a reference to another file by path.
type Inclusion interface {
	// GetLocation() returns content's original location.
	GetLocation() *Location

	// Path of the included content
	GetPath() string
}

// Context represents the current context of a Content being written.
type Context interface {
	// Get the value with the given name.
	Get(name string) (interface{}, bool)

	// Set the value for the given name.
	Set(name string, value interface{})

	// Remove the value with the given name.
	Remove(name string) interface{}

	// IsEmpty returns whether this context contains no values.
	IsEmpty() bool

	// ToStack converts this Context to a stack.
	ToStack() *ContextStack
}
