package mg

import (
	"fmt"
	"io"
)

type IncludeRawInstruction struct {
	UnscopedContent
	Text     string
	Path     string
	Origin   *Location
	Resolver FileResolver
}

var _ Inclusion = (*IncludeRawInstruction)(nil)
var _ Content = (*IncludeRawInstruction)(nil)

// NewIncludeRawInstruction creates a new IncludeRawInstruction.
func NewIncludeRawInstruction(arg string, location *Location, original string, resolver FileResolver) *IncludeRawInstruction {
	return &IncludeRawInstruction{Text: original, Path: arg, Origin: location, Resolver: resolver}
}

func (inc *IncludeRawInstruction) String() string {
	return fmt.Sprintf("IncludeRawInstruction{%s, %v, %v}", inc.Path, inc.Origin, inc.Resolver)
}

func (inc *IncludeRawInstruction) GetPath() string {
	return inc.Path
}

func (inc *IncludeRawInstruction) GetLocation() *Location {
	return inc.Origin
}

func (inc *IncludeRawInstruction) Write(writer io.Writer, context Context) ([]Content, error) {
	webFile, err := getInclusionByPath(inc, inc.Resolver, context, true)
	if err != nil {
		return nil, err
	}
	str, err := webFile.Processed.GetRawContents()
	if err != nil {
		return nil, err
	}
	content := NewStringContent(*str, webFile.GetLocation())
	return []Content{content}, nil
}
