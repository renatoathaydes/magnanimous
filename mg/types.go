package mg

import (
	"bytes"
	"errors"
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

type HasContent interface {
	AppendContent(content Content)
	Context() WebFileContext
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
	nestedStack  []HasContent
	NewExtension string
}

func (f *ProcessedFile) AppendContent(content Content) {
	s := len(f.nestedStack)
	if s > 0 {
		f.nestedStack[s-1].AppendContent(content)
	} else {
		f.Contents = append(f.Contents, content)
	}
	n, ok := content.(HasContent)
	if ok {
		f.nestedStack = append(f.nestedStack, n)
	}
}

func (f *ProcessedFile) EndNestedContent() error {
	s := len(f.nestedStack)
	if s > 0 {
		f.nestedStack = f.nestedStack[0 : s-1]
		return nil
	} else {
		return errors.New("end does not match any previous instruction")
	}
}

func (f *ProcessedFile) getFromNestedContent(name string) (interface{}, bool) {
	i := len(f.nestedStack) - 1
	for i > 0 {
		ctx := f.nestedStack[i].Context()
		if v, found := ctx[name]; found {
			return v, true
		}
	}
	return nil, false
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
