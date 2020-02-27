package inter

import (
	"bufio"
	"io"
	"unicode/utf8"
)

// SetWriter sets an alternative writer for output.
func (i *Interpreter) SetWriter(iow io.Writer) *Interpreter {
	i.writer = iow
	return i
}

// SetReader sets an alternative reader for input.
func (i *Interpreter) SetReader(ior io.Reader) *Interpreter {
	i.scanner = bufio.NewScanner(ior)
	//i.scanner.Split(bufio.ScanWords)
	i.scanner.Split(i.newSplitFunction())
	i.readingString = false // reset split function state
	return i
}

var _ bufio.SplitFunc

// newSplitFunction generates a split function dedicated to reading
// both tokens and string.
func (i *Interpreter) newSplitFunction() bufio.SplitFunc {
	i.readingString = false // state based on previous token
	return func(buf []byte, eof bool) (advance int, token []byte, err error) {
		if !i.readingString {
			advance, token, err = bufio.ScanWords(buf, eof)
			if string(token) == ".\"" {
				i.readingString = true
			}
			return advance, token, err
		}

		if i.readingString {

			// TODO - remove ending & finishing spaces

			start := 0
			for width, j := 0, start; j < len(buf); j += width {
				var r rune
				r, width = utf8.DecodeRune(buf[j:])
				if r == '"' {
					i.readingString = false
					return j + width, buf[start:j], nil
				}
			}
			// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.

			if eof && len(buf) > start {
				// switch back to normal mode
				i.readingString = false

				return len(buf), buf[start:], nil

			}

			return start, nil, nil // request more data
		}
		panic("invalid state")
	}
}
