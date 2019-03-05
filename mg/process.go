package mg

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (mag *Magnanimous) ReadAll() (WebFilesMap, error) {
	processedDir := filepath.Join(mag.SourcesDir, "processed")
	staticDir := filepath.Join(mag.SourcesDir, "static")

	procFiles, staticFiles, otherFiles := collectFiles(mag.SourcesDir, processedDir, staticDir)
	webFiles := WebFilesMap{
		WebFiles: make(map[string]WebFile, len(procFiles)+len(staticFiles)+len(otherFiles)),
	}

	err := ProcessAll(procFiles, processedDir, mag.SourcesDir, &webFiles)
	if err != nil {
		return webFiles, err
	}
	CopyAll(&staticFiles, staticDir, webFiles)
	AddNonWritables(&otherFiles, mag.SourcesDir, webFiles)

	return webFiles, nil
}

func ProcessAll(files []string, basePath, sourcesDir string, webFiles *WebFilesMap) error {
	resolver := DefaultFileResolver{BasePath: sourcesDir, Files: webFiles}
	for _, file := range files {
		wf, err := ProcessFile(file, basePath, &resolver)
		if err != nil {
			return err
		}
		webFiles.WebFiles[file] = *wf
	}
	return nil
}

func ProcessFile(file, basePath string, resolver FileResolver) (*WebFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, &MagnanimousError{message: err.Error(), Code: IOError}
	}
	reader := bufio.NewReader(f)
	s, err := f.Stat()
	if err != nil {
		return nil, &MagnanimousError{message: err.Error(), Code: IOError}
	}
	processed, magErr := ProcessReader(reader, file, int(s.Size()), resolver)
	if magErr != nil {
		return nil, magErr
	}
	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: processed, NonWritable: nonWritable}, nil
}

func ProcessReader(reader *bufio.Reader, file string, sizeHint int, resolver FileResolver) (*ProcessedFile, error) {
	var builder strings.Builder
	builder.Grow(sizeHint)
	isMarkDown := isMd(file)
	processed := ProcessedFile{Path: file}
	stack := []ContentContainer{&processed}
	state := parserState{file: file, row: 1, col: 1, builder: &builder, reader: reader, contentStack: stack}
	magErr := parseText(&state, resolver)
	if magErr != nil {
		return &processed, magErr
	}
	if isMarkDown {
		processed = MarkdownToHtml(processed)
	}
	return &processed, nil
}

func WriteTo(dir string, filesMap WebFilesMap) error {
	stack := ContextStack{}
	if globalCtx, ok := filesMap.WebFiles["source/_global_context"]; ok {
		ctx := globalCtx.Processed.ResolveContext(filesMap, stack)
		stack = NewContextStack(ctx)
	}

	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	for file, wf := range filesMap.WebFiles {
		if wf.NonWritable {
			continue
		}
		targetPath, err := filepath.Rel(wf.BasePath, file)
		if err != nil {
			log.Printf("Unable to relativize path %s", file)
			targetPath = file
		}
		targetFile := filepath.Join(dir, targetPath)
		if wf.Processed.NewExtension != "" {
			ext := filepath.Ext(targetFile)
			targetFile = targetFile[0:len(targetFile)-len(ext)] + wf.Processed.NewExtension
		}
		magErr := writeFile(file, targetFile, wf, filesMap, stack)
		if magErr != nil {
			return magErr
		}
	}
	return nil
}

func writeFile(file, targetFile string, wf WebFile, filesMap WebFilesMap, stack ContextStack) error {
	log.Printf("Creating file %s from %s", targetFile, file)
	err := os.MkdirAll(filepath.Dir(targetFile), 0770)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	f, err := os.Create(targetFile)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, c := range wf.Processed.contents {
		err := c.Write(w, filesMap, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wf *WebFile) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	return writeContents(wf.Processed, writer, files, stack)
}

func writeContents(cc ContentContainer, writer io.Writer, files WebFilesMap,
	stack ContextStack) error {
	for _, c := range cc.GetContents() {
		err := c.Write(writer, files, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func inclusionChainToString(locations []ContextStackItem) string {
	var b strings.Builder
	b.WriteRune('[')
	last := len(locations) - 1
	for i, loc := range locations {
		b.WriteString(loc.Location.String())
		if i != last {
			b.WriteString(" -> ")
		}
	}
	b.WriteRune(']')
	return b.String()
}
