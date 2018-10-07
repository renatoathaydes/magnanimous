package mg

import "io"

type WebFilesMap map[string]WebFile

type WebFileContext map[string]interface{}

type WebFile struct {
	BasePath  string
	Context   *WebFileContext
	Processed *ProcessedFile
}

type Content interface {
	Write(writer io.Writer)
}

type ProcessedFile struct {
	Contents []Content
}

func (f *ProcessedFile) appendContent(content Content) {
	f.Contents = append(f.Contents, content)
}
