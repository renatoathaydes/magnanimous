package mg

import (
	"fmt"
	"io"
)

type StringContent struct {
	Text string
}

func (c *StringContent) Write(writer io.Writer, stack ContextStack) error {
	_, err := writer.Write([]byte(c.Text))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (c *StringContent) String() string {
	return fmt.Sprintf("StringContent{%s}", c.Text)
}
