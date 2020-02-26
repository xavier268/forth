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
	i.scanner.Split(newSplitFunction())
	return i
}

var _ bufio.SplitFunc

// newSplitFunction generates a split function dedicated to reading
// both tokens and string.
func newSplitFunction() bufio.SplitFunc {
	var readingString = false // state based on previous token
	return func(buf []byte, eof bool) (advance int, token []byte, err error) {
		if !readingString {
			advance, token, err = bufio.ScanWords(buf, eof)
			if string(token) == ".\"" {
				readingString = true
			}
			return advance, token, err
		}

		if readingString { // TODO Check vs golang package implementation ...
			start := 0
			for width, i := 0, start; i < len(buf); i += width {
				var r rune
				r, width = utf8.DecodeRune(buf[i:])
				if r == '"' {
					readingString = false
					return i + width, buf[start:i], nil
				}

			}

			return 0, nil, ErrUnexpectedEndOfLine
		}
		panic("invalid state")
	}
}
