package mg

import (
	"fmt"
	"io"
)

type IncludeInstruction struct {
	UnscopedContent
	Text     string
	Path     string
	Origin   *Location
	Resolver FileResolver
}

var _ Inclusion = (*IncludeInstruction)(nil)
var _ Content = (*IncludeInstruction)(nil)

func NewIncludeInstruction(arg string, location *Location, original string, resolver FileResolver) *IncludeInstruction {
	return &IncludeInstruction{Text: original, Path: arg, Origin: location, Resolver: resolver}
}

func (inc *IncludeInstruction) String() string {
	return fmt.Sprintf("IncludeInstruction{%s, %v, %v}", inc.Path, inc.Origin, inc.Resolver)
}

func (inc *IncludeInstruction) GetPath() string {
	return inc.Path
}

func (inc *IncludeInstruction) GetLocation() *Location {
	return inc.Origin
}

func (inc *IncludeInstruction) Write(writer io.Writer, context Context) ([]Content, error) {
	webFile, err := getInclusionByPath(inc, inc.Resolver, context, true)
	if err != nil {
		return nil, err
	}
	return webFile.Processed.GetContents(), nil
}
