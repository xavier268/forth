package inter

import (
	"bufio"
	"io"
)

// SetWriter sets an alternative writer for output.
func (i *Interpreter) SetWriter(iow io.Writer) *Interpreter {
	i.writer = iow
	return i
}

// SetReader sets an alternative reader for input.
func (i *Interpreter) SetReader(ior io.Reader) *Interpreter {
	i.scanner = bufio.NewScanner(ior)
	i.scanner.Split(bufio.ScanWords)
	return i
}
