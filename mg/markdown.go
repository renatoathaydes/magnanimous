package mg

import (
	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
	"io"
)

type HtmlFromMarkdownContent struct {
	MarkDownContent []Content
}

var _ Content = (*HtmlFromMarkdownContent)(nil)

var chromaRenderer = blackfriday.WithRenderer(bfchroma.NewRenderer(
	bfchroma.WithoutAutodetect(), bfchroma.Style("lovelace")))

func MarkdownToHtml(file ProcessedFile) ProcessedFile {
	return ProcessedFile{
		Path:         file.Path,
		contents:     []Content{&HtmlFromMarkdownContent{MarkDownContent: file.contents}},
		NewExtension: ".html",
	}
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	mdBytes, err := asBytes(f.MarkDownContent, files, stack)

	md := blackfriday.Run(mdBytes, chromaRenderer)

	_, err = writer.Write(md)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}
