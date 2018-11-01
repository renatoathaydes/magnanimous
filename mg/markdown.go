package mg

import (
	"bytes"
	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
	"io"
)

var chromaRenderer = blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.WithoutAutodetect()))

func MarkdownToHtml(file ProcessedFile) ProcessedFile {
	convertedContent := make([]Content, 0, len(file.contents))
	for _, c := range file.contents {
		if c.IsMarkDown() {
			convertedContent = append(convertedContent, &HtmlFromMarkdownContent{MarkDownContent: c})
		} else {
			convertedContent = append(convertedContent, c)
		}
	}
	return ProcessedFile{contents: convertedContent, context: file.context, NewExtension: ".html"}
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	content, magErr := readBytes(&f.MarkDownContent, files, inclusionChain)
	if magErr != nil {
		return magErr
	}
	_, err := writer.Write(blackfriday.Run(content, chromaRenderer))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (_ *HtmlFromMarkdownContent) IsMarkDown() bool {
	return false
}

func readBytes(c *Content, files WebFilesMap, inclusionChain []InclusionChainItem) ([]byte, error) {
	var b bytes.Buffer
	b.Grow(1024)
	err := (*c).Write(&b, files, inclusionChain)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
