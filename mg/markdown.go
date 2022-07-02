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

func isInMd(stack ContextStack) bool {
	s := stack.locations
	len := len(s)
	if len > 0 {
		currentIsMd := isMd(s[len-1].Origin)
		if len < 2 || !currentIsMd {
			return currentIsMd
		}
		// if MD is included from non-MD, we say we're not in MD
		return isMd(s[len-2].Origin)
	}
	return false
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
