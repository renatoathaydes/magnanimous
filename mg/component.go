package mg

import (
	"io"
	"log"
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

func (c *Component) AppendContent(content Content) {
	c.contents = append(c.GetContents(), content)
}

func (c *Component) GetContents() []Content {
	if isMd(c.Location.Origin) {
		return []Content{&HtmlFromMarkdownContent{MarkDownContent: c.contents}}
	} else {
		return c.contents
	}
}

func NewComponentInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	return &Component{
		Path:     arg,
		Location: location,
		Text:     original,
		resolver: resolver,
	}
}

func (c *Component) Write(writer io.Writer, stack ContextStack) error {
	params := magParams{stack: stack, fileResolver: c.resolver, location: c.Location}
	maybePath := pathOrEval(c.Path, &params)
	var actualPath string
	if s, ok := maybePath.(string); ok {
		actualPath = s
	} else {
		log.Printf("WARNING: path expression evaluated to invalid value: %v", maybePath)
		_, err := writer.Write([]byte(c.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
		return nil
	}
	componentFile, ok := params.File(actualPath)
	if !ok {
		log.Printf("WARNING: (%s) refers to a non-existent Component: %s", c.Location, actualPath)
		_, err := writer.Write([]byte(c.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		stack = stack.Push(c.Location, true)
		err := detectCycle(stack, actualPath, componentFile.Processed.Path, c.Location)
		if err != nil {
			return err
		}
		contents, err := body(c, stack)
		if err != nil {
			return err
		}
		stack.Top().Set("__contents__", string(contents))
		err = componentFile.Write(writer, stack)
		if err != nil {
			return err
		}
	}
	return nil
}
