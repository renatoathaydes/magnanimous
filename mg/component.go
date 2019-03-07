package mg

import (
	"io"
	"log"
)

type Component struct {
	Path     string
	Location *Location
	Text     string
	Resolver FileResolver
	contents []Content
}

var _ Content = (*Component)(nil)
var _ ContentContainer = (*Component)(nil)

func (c *Component) AppendContent(content Content) {
	c.contents = append(c.contents, content)
}

func (c *Component) GetContents() []Content {
	return c.contents
}

func NewComponentInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	return &Component{
		Path:     arg,
		Location: location,
		Text:     original,
		Resolver: resolver,
	}
}

func (c *Component) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	path := c.Resolver.Resolve(c.Path, c.Location)
	//fmt.Printf("Including %s from %v : %s\n", c.Path, c.Origin, path)
	componentFile, ok := files.WebFiles[path]
	if !ok {
		log.Printf("WARNING: (%s) refers to a non-existent Component: %s", c.Location, c.Path)
		_, err := writer.Write([]byte(c.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		stack = stack.Push(c.Location)
		err := detectCycle(stack, c.Path, path, c.Location)
		if err != nil {
			return err
		}
		contents, err := body(c, files, stack)
		if err != nil {
			return err
		}
		stack.Top().Context.Set("__contents__", string(contents))
		err = componentFile.Write(writer, files, stack)
		if err != nil {
			return err
		}
	}
	return nil
}
