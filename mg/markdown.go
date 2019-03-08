package mg

import (
	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
	"io"
)

// HtmlFromMarkdownContent wraps its contents so that it can convert that to HTML on write.
//
// Any included files are written without conversion. Hence, when a HTML file, for example, in included from a
// markdown file, its contents will be written as they are rather than passed through the markdown-to-html engine.
type HtmlFromMarkdownContent struct {
	MarkDownContent []Content
}

var _ Content = (*HtmlFromMarkdownContent)(nil)

var mdStyle = bfchroma.Style("lovelace")

// MarkdownToHtml wraps the contents of a [ProcessedFile] so that it will be converted from markdown
// to HTML on write.
func MarkdownToHtml(file ProcessedFile) ProcessedFile {
	return ProcessedFile{
		Path:         file.Path,
		contents:     []Content{&HtmlFromMarkdownContent{MarkDownContent: file.contents}},
		NewExtension: ".html",
	}
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	// accumulate all contents that do not include another file, then writeAsHtmlAndReset it...
	// inclusions are written as they are, without translation from markdown to html.
	var nonIncludedContent []Content
	var err error

	for _, c := range f.MarkDownContent {
		switch c.(type) {
		case *IncludeInstruction:
			nonIncludedContent, err = writeAsHtmlAndReset(nonIncludedContent, writer, c, files, stack)
			if err != nil {
				return err
			}
		case *Component:
			nonIncludedContent, err = writeAsHtmlAndReset(nonIncludedContent, writer, c, files, stack)
			if err != nil {
				return err
			}
		default:
			nonIncludedContent = append(nonIncludedContent, c)
		}
	}

	err = writeAsHtml(nonIncludedContent, writer, files, stack)
	return err
}

func unwrapMarkdownContent(f *ProcessedFile) ([]Content, bool) {
	if len(f.contents) == 1 {
		if md, ok := f.contents[0].(*HtmlFromMarkdownContent); ok {
			return md.MarkDownContent, true
		}
	}
	return nil, false
}

func writeAsHtmlAndReset(contents []Content, writer io.Writer, content Content, files WebFilesMap,
	stack ContextStack) ([]Content, error) {
	// write contents, converting it to html
	err := writeAsHtml(contents, writer, files, stack)
	if err != nil {
		return nil, err
	}
	// write content as it is, returning nil (or empty array)
	return nil, content.Write(writer, files, stack)
}

func writeAsHtml(c []Content, writer io.Writer, files WebFilesMap, stack ContextStack) error {
	if len(c) == 0 {
		return nil
	}
	mdBytes, err := asBytes(c, files, stack)

	var chromaRenderer = blackfriday.WithRenderer(
		bfchroma.NewRenderer(bfchroma.WithoutAutodetect(), mdStyle))

	md := blackfriday.Run(mdBytes, chromaRenderer)

	_, err = writer.Write(md)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}
