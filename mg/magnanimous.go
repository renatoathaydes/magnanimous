package mg

import (
	"bufio"
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

	err := ProcessAll(procFiles, processedDir, mag.SourcesDir, &webFiles)
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

// ProcessFile processes the given file.
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
	processed, magErr := ProcessReader(reader, file, int(s.Size()), resolver, s.ModTime())
	if magErr != nil {
		return nil, magErr
	}

	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: processed, NonWritable: nonWritable}, nil
}

// ProcessReader processes the contents provided by the given reader.
func ProcessReader(reader *bufio.Reader, file string, sizeHint int, resolver FileResolver,
	lastUpdated time.Time) (*ProcessedFile, error) {

	var builder strings.Builder
	builder.Grow(sizeHint)
	isMarkDown := isMd(file)
	processed := ProcessedFile{Path: file, LastUpdated: lastUpdated}
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
		ctx := globalCtx.Processed.ResolveContext(stack)
		stack = NewContextStack(ctx)
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
		magErr := writeFile(file, targetFile, wf, stack)
		if magErr != nil {
			return magErr
		}
	}
	return nil
}

func writeFile(file, targetFile string, wf WebFile, stack ContextStack) error {
	log.Printf("Creating file %s from %s", targetFile, file)
	stack = stack.Push(nil, true)
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
		err := c.Write(w, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wf *WebFile) Write(writer io.Writer, stack ContextStack) error {
	return writeContents(wf.Processed, writer, stack)
}

func writeContents(cc ContentContainer, writer io.Writer, stack ContextStack) error {
	for _, c := range cc.GetContents() {
		err := c.Write(writer, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func inclusionChainToString(inclusionChain []Location) string {
	var b strings.Builder
	b.WriteRune('[')
	var includes []string
	for _, loc := range inclusionChain {
		includes = append(includes, loc.String())
	}
	b.WriteString(strings.Join(includes, " -> "))
	b.WriteRune(']')
	return b.String()
}
