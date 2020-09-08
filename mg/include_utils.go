package mg

import (
	"fmt"
	"log"
	"strings"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

func pathOrEval(path string, params *magParams) interface{} {
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
		res, err := expression.Eval(path[startIndex:], params)
		if err != nil {
			log.Printf("WARNING: eval expression error: %v", err)
		} else {
			return res
		}
	}
	return path
}

func detectCycle(stack ContextStack, includedPath, absPath string, location *Location) error {
	for _, loc := range stack.locations {
		if loc.Origin == absPath {
			chain := inclusionChainToString(stack.locations)
			return &MagnanimousError{
				Code: InclusionCycleError,
				message: fmt.Sprintf(
					"Cycle detected! Inclusion of %s at %s comes back into itself via %s",
					includedPath, location.String(), chain),
			}
		}
	}
	return nil
}
