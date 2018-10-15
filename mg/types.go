package mg

import (
	"bytes"
	"fmt"
	"io"
)

type WebFilesMap map[string]WebFile

type WebFileContext map[string]interface{}

type WebFile struct {
	BasePath    string
	Context     WebFileContext
	Processed   ProcessedFile
	NonWritable bool
}

type Location struct {
	Origin string
	Row    uint32
	Col    uint32
}

type Content interface {
	Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError
	IsMarkDown() bool
}

type StringContent struct {
	Text     string
	MarkDown bool
}

type HtmlFromMarkdownContent struct {
	MarkDownContent Content
}

type ProcessedFile struct {
	Contents     []Content
	NewExtension string
}

func (f *ProcessedFile) AppendContent(content Content) {
	f.Contents = append(f.Contents, content)
}

func (f *ProcessedFile) Bytes(files WebFilesMap, inclusionChain []Location) []byte {
	var b bytes.Buffer
	b.Grow(512)
	for _, c := range f.Contents {
		if c != nil {
			c.Write(&b, files, inclusionChain)
		}
	}
	return b.Bytes()
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Origin, l.Row, l.Col)
}
