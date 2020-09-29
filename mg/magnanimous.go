package mg

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// ReadAll source files, creating a mapping from file paths to [WebFile] instances.
func (mag *Magnanimous) ReadAll() (WebFilesMap, error) {
	processedDir := filepath.Join(mag.SourcesDir, "processed")
	staticDir := filepath.Join(mag.SourcesDir, "static")

	procFiles, staticFiles, otherFiles := collectFiles(mag.SourcesDir, processedDir, staticDir)
	webFiles := WebFilesMap{
		WebFiles: make(map[string]WebFile, len(procFiles)+len(staticFiles)+len(otherFiles)),
	}

	err := mag.ProcessAll(procFiles, processedDir, &webFiles)
	if err != nil {
		return webFiles, err
	}
	err = CopyAll(&staticFiles, staticDir, webFiles)
	if err != nil {
		return webFiles, err
	}
	err = AddNonWritables(&otherFiles, mag.SourcesDir, webFiles)
	return webFiles, err
}

// ProcessAll given files, putting the results in the given webFiles map.
func (mag *Magnanimous) ProcessAll(files []string, basePath string, webFiles *WebFilesMap) error {
	resolver := DefaultFileResolver{BasePath: mag.SourcesDir, Files: webFiles}
	for _, file := range files {
		wf, err := ProcessFile(file, basePath, &resolver)
		if err != nil {
			return err
		}
		webFiles.WebFiles[file] = *wf
	}
	return nil
}

// ProcessFile processes the given file.
func ProcessFile(file, basePath string, resolver FileResolver) (*WebFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	processed, err := ProcessReader(reader, file, basePath, int(s.Size()), resolver, s.ModTime())
	if err != nil {
		return nil, err
	}

	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: processed, NonWritable: nonWritable}, nil
}

// ProcessReader processes the contents provided by the given reader.
func ProcessReader(reader *bufio.Reader, file, basePath string, sizeHint int, resolver FileResolver,
	lastUpdated time.Time) (*ProcessedFile, error) {

	var builder strings.Builder
	builder.Grow(sizeHint)
	processed := ProcessedFile{BasePath: basePath, Path: file, LastUpdated: lastUpdated}
	if isMd(file) {
		processed.NewExtension = "html"
	}
	stack := []ContentContainer{&processed}
	state := parserState{file: file, row: 1, col: 1, builder: &builder, reader: reader, contentStack: stack}
	magErr := parseText(&state, resolver)
	if magErr != nil {
		return &processed, magErr
	}
	return &processed, nil
}

func (mag *Magnanimous) newContextStack(filesMap WebFilesMap) ContextStack {
	var stack = NewContextStack(NewContext())
	var globalCtxPath string
	if mag.GlobalContex != "" {
		globalCtxPath = path.Join(mag.SourcesDir, "processed", mag.GlobalContex)
	} else {
		globalCtxPath = path.Join(mag.SourcesDir, "processed", "_global_context")
	}
	if globalCtx, ok := filesMap.WebFiles[globalCtxPath]; ok {
		log.Printf("Using global context file: %s", globalCtxPath)
		globalCtx.Processed.ResolveContext(&stack, true)
	} else if mag.GlobalContex != "" {
		log.Printf("WARNING: global context file was not found: %s", globalCtxPath)
	} else {
		log.Println("No global context file defined.")
	}
	return stack
}

// WriteTo writes all files in the given map on the given directory.
func (mag *Magnanimous) WriteTo(dir string, filesMap WebFilesMap) error {
	stack := mag.newContextStack(filesMap)

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
			targetFile = changeFileExt(targetFile, wf.Processed.NewExtension)
		}
		magErr := writeFile(file, targetFile, wf, &stack)
		if magErr != nil {
			return magErr
		}
	}
	return nil
}

func writeFile(file, targetFile string, wf WebFile, stack *ContextStack) error {
	err := os.MkdirAll(filepath.Dir(targetFile), 0770)
	if err != nil {
		return err
	}

	if wf.SkipIfUpToDate {
		upToDate, err := isUpToDate(&wf, targetFile)
		if err != nil {
			return err
		}
		if upToDate {
			log.Printf("Skipping file %s as it has not been updated since last run.", targetFile)
			return nil
		}
	}

	log.Printf("Creating file %s from %s", targetFile, file)
	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	return wf.Write(w, stack, true, false)
}

func (wf *WebFile) Write(writer io.Writer, stack *ContextStack, useScope, writePlain bool) error {
	if useScope {
		pushResult := stack.Push(wf.GetLocation(), true)
		defer stack.Pop(pushResult)
	}

	inMd := isMd(wf.Processed.GetLocation().Origin)
	var buffer bytes.Buffer
	buffer.Grow(512)
	inMd, err := writeContents(wf.Processed.GetContents(), writer, &buffer, stack, inMd, writePlain)
	if err == nil {
		if inMd && !writePlain {
			err = flushMdAsHtml(&buffer, writer)
		} else {
			err = flush(&buffer, writer)
		}
	}
	return err
}

func (wf *WebFile) GetLocation() *Location {
	origin := filepath.Join(wf.BasePath, wf.Name)
	return &Location{Origin: origin, Row: 0, Col: 0}
}

func writeContents(contents []Content, writer io.Writer, buffer *bytes.Buffer, stack *ContextStack,
	inMd, writePlain bool) (stillInMd bool, err error) {
	stillInMd = inMd
	for _, content := range contents {
		stillInMd, err = writeContent(content, writer, buffer, stack, stillInMd, writePlain)
		if err != nil {
			return
		}
	}
	return
}

func writeContent(c Content, writer io.Writer, buffer *bytes.Buffer, stack *ContextStack,
	inMd, writePlain bool) (stillInMd bool, err error) {
	isScoped := c.IsScoped()
	stillInMd = isMd(c.GetLocation().Origin)
	// flush the buffer only in case we're writing plain content, or the content switched format from/to md.
	if writePlain || (!inMd && stillInMd) {
		err = flush(buffer, writer)
	} else if inMd && !stillInMd {
		err = flushMdAsHtml(buffer, writer)
	}
	if err != nil {
		return
	}

	pushResult := stack.Push(c.GetLocation(), isScoped)
	defer stack.Pop(pushResult)

	next, err := c.Write(buffer, stack)
	if err == nil && len(next) > 0 {
		stillInMd, err = writeContents(next, writer, buffer, stack, stillInMd, writePlain)
	}
	return
}

func flush(buffer *bytes.Buffer, writer io.Writer) error {
	defer buffer.Reset()
	bytes := buffer.Bytes()
	if len(bytes) == 0 {
		return nil
	}
	_, err := writer.Write(bytes)
	return err
}

func isUpToDate(wf *WebFile, targetFile string) (bool, error) {
	stat, err := os.Stat(targetFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return stat.ModTime().After(wf.Processed.LastUpdated), nil
}
