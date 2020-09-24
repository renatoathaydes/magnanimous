package mg

import (
	"io"
)

type Component struct {
	Path     string
	Location *Location
	Text     string
	resolver FileResolver
	contents []Content
}

var _ Content = (*Component)(nil)
var _ ContentContainer = (*Component)(nil)
var _ Inclusion = (*Component)(nil)

type internalComponent struct {
	UnscopedContent
	location *Location
	contents []Content
}

var _ Content = (*internalComponent)(nil)

func NewComponentInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	return &Component{
		Path:     arg,
		Location: location,
		Text:     original,
		resolver: resolver,
	}
}

func (c *Component) AppendContent(content Content) {
	c.contents = append(c.contents, content)
}

func (c *Component) GetPath() string {
	return c.Path
}

func (c *Component) GetLocation() *Location {
	return c.Location
}

func (c *Component) IsScoped() bool {
	return true
}

func (c *Component) Write(writer io.Writer, context Context) ([]Content, error) {
	componentFile, err := getInclusionByPath(c, c.resolver, context, true)
	if err != nil {
		return nil, err
	}

	resolveContext(c.contents, context)

	context.Set("__contents__", &internalComponent{location: c.Location, contents: c.contents})

	return componentFile.Processed.GetContents(), nil
}

func (c *internalComponent) GetLocation() *Location {
	return c.location
}

func (c *internalComponent) Write(writer io.Writer, context Context) ([]Content, error) {
	return c.contents, nil
}
