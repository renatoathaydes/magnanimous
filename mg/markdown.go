package mg

import (
	"io"

	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
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

// SetCodeStyle sets the code style used to highlight source code.
// See https://xyproto.github.io/splash/docs/all.html for the supported styles.
func SetCodeStyle(style string) {
	mdStyle = bfchroma.Style(style)
}

// MarkdownToHtml wraps the contents of a [ProcessedFile] so that it will be converted from markdown
// to HTML on write.
func MarkdownToHtml(file ProcessedFile) ProcessedFile {
	return ProcessedFile{
		Path:         file.Path,
		LastUpdated:  file.LastUpdated,
		contents:     []Content{&HtmlFromMarkdownContent{MarkDownContent: file.contents}},
		NewExtension: "html",
	}
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, stack ContextStack) error {
	// accumulate all markdown content until pure HTML content is found, then writeAsHtmlAndReset it...
	var markdownContent []Content
	var err error

	for _, c := range f.MarkDownContent {
		switch v := c.(type) {
		case *IncludeInstruction:
			if v.AsHTML {
				markdownContent, err = writeAsHtmlAndReset(markdownContent, writer, c, stack)
				if err != nil {
					return err
				}
			} else {
				markdownContent = append(markdownContent, c)
			}
		case *Component:
			if v.AsHTML {
				markdownContent, err = writeAsHtmlAndReset(markdownContent, writer, c, stack)
				if err != nil {
					return err
				}
			} else {
				markdownContent = append(markdownContent, c)
			}
		default:
			markdownContent = append(markdownContent, c)
		}
	}

	err = writeAsHtml(markdownContent, writer, stack)
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

func writeAsHtmlAndReset(markdownContents []Content, writer io.Writer, htmlContent Content,
	stack ContextStack) ([]Content, error) {
	// write markdownContents, converting it to html
	err := writeAsHtml(markdownContents, writer, stack)
	if err != nil {
		return nil, err
	}
	// write htmlContent as it is, returning nil (or empty array) to "reset" the markdown contents slice
	return nil, htmlContent.Write(writer, stack)
}

func writeAsHtml(c []Content, writer io.Writer, stack ContextStack) error {
	if len(c) == 0 {
		return nil
	}
	mdBytes, err := asBytes(c, stack)

	var chromaRenderer = blackfriday.WithRenderer(
		bfchroma.NewRenderer(bfchroma.WithoutAutodetect(), mdStyle))

	html := blackfriday.Run(mdBytes, chromaRenderer)

	_, err = writer.Write(html)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}
