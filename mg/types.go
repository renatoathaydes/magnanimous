package mg

import (
	"fmt"
	"io"
)

type WebFilesMap map[string]WebFile

type WebFileContext map[string]interface{}

type WebFile struct {
	BasePath    string
	Context     *WebFileContext
	Processed   *ProcessedFile
	NonWritable bool
}

type Location struct {
	Origin string
	Row    uint32
	Col    uint32
}

type Content interface {
	Write(writer io.Writer, files WebFilesMap)
}

type StringContent struct {
	Text string
}

type IncludeInstruction struct {
	Name   string
	Path   string
	Origin Location
}

type ProcessedFile struct {
	Contents []Content
}

func (f *ProcessedFile) AppendContent(content Content) {
	f.Contents = append(f.Contents, content)
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Origin, l.Row, l.Col)
}
