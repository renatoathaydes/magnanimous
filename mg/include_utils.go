package mg

import (
	"fmt"
	"log"
	"strings"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

func getInclusionByPath(inc Inclusion, resolver FileResolver, context Context, checkCycles bool) (*WebFile, error) {
	maybePath := pathOrEval(inc.GetPath(), context)
	var actualPath string
	if s, ok := maybePath.(string); ok {
		actualPath = s
	} else {
		return nil, fmt.Errorf("path expression evaluated to non-string value: %v", maybePath)
	}
	f := resolver.Resolve(actualPath, inc.GetLocation(), context.ToStack().NearestLocation())
	webFile, ok := resolver.Get(f)
	if !ok {
		return nil, fmt.Errorf("path expression refers non-existent resource: %s", actualPath)
	}
	if checkCycles {
		err := detectCycle(context, actualPath, webFile.Processed.Path, inc.GetLocation())
		if err != nil {
			return nil, err
		}
	}
	return webFile, nil
}

func pathOrEval(path string, context Context) interface{} {
	startIndex := -1
	if strings.HasPrefix(path, "eval ") {
		startIndex = 5
	}
	if strings.HasPrefix(path, "\"") ||
		strings.HasPrefix(path, "`") ||
		strings.HasPrefix(path, "[") {
		startIndex = 0
	}
	if startIndex != -1 {
		// treat rest of argument as an expression that evaluates to a path
		res, err := expression.Eval(path[startIndex:], context)
		if err != nil {
			log.Printf("WARNING: eval expression error: %v", err)
		} else {
			return res
		}
	}
	return path
}

func detectCycle(context Context, includedPath, absPath string, location *Location) error {
	for _, loc := range context.ToStack().locations {
		if loc.Origin == absPath {
			chain := inclusionChainToString(context.ToStack().locations)
			return fmt.Errorf("Cycle detected! Inclusion of %s at %s comes back into itself via %s",
				includedPath, location.String(), chain)
		}
	}
	return nil
}

func inclusionChainToString(inclusionChain []*Location) string {
	// the first location is expected to be the name of the file being written
	if len(inclusionChain) > 1 {
		inclusionChain = inclusionChain[1:len(inclusionChain)]
	}
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
