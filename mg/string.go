package mg

import (
	"fmt"
	"io"
)

// StringContent is a simple Writable Content wrapping a [string].
type StringContent struct {
	UnscopedContent
	text     string
	location *Location
}

var _ Content = (*StringContent)(nil)

func NewStringContent(text string, location *Location) *StringContent {
	return &StringContent{text: text, location: location}
}

func unevaluatedExpression(original string, location *Location) Content {
	return NewStringContent(fmt.Sprintf("{{%s}}", original), location)
}

func unevaluatedExpressions(original string, location *Location) []Content {
	return []Content{unevaluatedExpression(original, location)}
}

func (c *StringContent) Write(writer io.Writer, context Context) ([]Content, error) {
	_, err := writer.Write([]byte(c.text))
	return nil, err
}

func (c *StringContent) GetLocation() *Location {
	return c.location
}

func (c *StringContent) String() string {
	return fmt.Sprintf("StringContent{%s}", c.text)
}
