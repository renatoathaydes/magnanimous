package mg

import (
	"bytes"
	"io"

	"github.com/Depado/bfchroma/v2"
	"github.com/russross/blackfriday/v2"
)

var mdStyle = bfchroma.Style("lovelace")

// SetCodeStyle sets the code style used to highlight source code.
// See https://xyproto.github.io/splash/docs/all.html for the supported styles.
func SetCodeStyle(style string) {
	mdStyle = bfchroma.Style(style)
}

func flushMdAsHtml(buffer *bytes.Buffer, writer io.Writer) error {
	defer buffer.Reset()
	mdBytes := buffer.Bytes()
	if len(mdBytes) == 0 {
		return nil
	}
	var chromaRenderer = blackfriday.WithRenderer(
		bfchroma.NewRenderer(bfchroma.WithoutAutodetect(), mdStyle))

	html := blackfriday.Run(mdBytes, chromaRenderer)

	_, err := writer.Write(html)
	return err
}
