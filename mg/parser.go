package mg

import (
	"fmt"
	"io"
)

func parseText(state *parserState, resolver FileResolver) error {
	previousWasOpenBracket := false
	previousWasEscape := false
	reader := state.reader
	builder := state.builder

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &MagnanimousError{message: err.Error(), Code: IOError}
		}

		if r == '\n' {
			state.row++
			state.col = 1
			if !previousWasEscape {
				builder.WriteRune(r)
			} else {
				previousWasEscape = false
			}
			continue
		}

		if r == '\r' {
			if !previousWasEscape {
				builder.WriteRune(r)
			}
			continue
		}

		state.col++

		if r == '{' {
			if previousWasEscape {
				builder.WriteRune(r)
				previousWasEscape = false
				continue
			}
			if previousWasOpenBracket {
				if builder.Len() > 0 {
					state.pf.AppendContent(&StringContent{Text: builder.String()})
					builder.Reset()
				}
				previousWasOpenBracket = false
				magErr := parseInstruction(state, resolver)
				if magErr != nil {
					return magErr
				}
			} else {
				previousWasOpenBracket = true
			}
			continue
		}

		if previousWasOpenBracket {
			builder.WriteRune('{')
		}
		previousWasOpenBracket = false

		if r == '\\' {
			previousWasEscape = true
		} else {
			builder.WriteRune(r)
		}
	}

	// append any pending content to the builder
	if builder.Len() > 0 || previousWasOpenBracket || previousWasEscape {
		if previousWasOpenBracket {
			builder.WriteRune('{')
		} else if previousWasEscape {
			builder.WriteRune('\\')
		}
		state.pf.AppendContent(&StringContent{Text: builder.String()})
		builder.Reset()
	}

	return nil
}

func parseInstruction(state *parserState, resolver FileResolver) error {
	previousWasCloseBracket := false
	previousWasEscape := false
	instrFirstRow := state.row
	instrFirstCol := state.col - 2
	reader := state.reader
	builder := state.builder

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &MagnanimousError{message: err.Error(), Code: IOError}
		}

		if r == '\n' {
			state.row++
			state.col = 1
			builder.WriteRune(r)
			continue
		}

		state.col++

		if r == '}' {
			if previousWasEscape {
				builder.WriteRune(r)
				previousWasEscape = false
				continue
			}
			if previousWasCloseBracket {
				if builder.Len() > 0 {
					appendContent(state.pf, builder.String(),
						Location{Origin: state.file, Row: instrFirstRow, Col: instrFirstCol},
						resolver)
				}
				builder.Reset()
				return nil
			} else {
				previousWasCloseBracket = true
			}
			continue
		}

		if previousWasCloseBracket {
			builder.WriteRune('}')
		}
		previousWasCloseBracket = false

		if r == '\\' {
			previousWasEscape = true
		} else {
			builder.WriteRune(r)
		}
	}

	return NewError(Location{Origin: state.file, Row: state.row, Col: state.col}, ParseError,
		fmt.Sprintf("instruction started at (%d:%d) was not properly closed with '}}'",
			instrFirstRow, instrFirstCol))
}
