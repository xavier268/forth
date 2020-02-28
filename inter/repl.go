package inter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Color definitions.
const (
	ColorRed   = "\033[0;31m"
	ColorGreen = "\033[1;32m"
	ColorBlue  = "\033[1;34m"

	ColorOff = "\033[m"
)

// Prompt the user for entry
func (i *Interpreter) Prompt() string {
	// i.dump()
	if i.compileMode {
		return fmt.Sprintf("\n%scompile:%s ", ColorBlue, ColorOff)
	}
	pt := fmt.Sprintf("\n%s", ColorGreen)
	if len(i.rs.data) > 4 {
		pt += fmt.Sprintf("rs(top 4)%v", i.rs.data[len(i.rs.data)-4:])
	} else {
		pt += fmt.Sprintf("rs%v", i.rs.data)
	}
	if len(i.ds.data) > 4 {
		pt += fmt.Sprintf(", ds(top 4)%v", i.ds.data[len(i.ds.data)-4:])
	} else {
		pt += fmt.Sprintf(", ds%v", i.ds.data)
	}
	pt += fmt.Sprintf("%s> ", ColorOff)

	return pt
}

// Repl is the main Read-Evaluate-Print-Loop.
// Does not use the scanner set in Interpreter.
func (i *Interpreter) Repl() {

	// scanner reads line per line
	linescan := bufio.NewScanner(os.Stdin)
	linescan.Split(bufio.ScanLines)

	for { // repl loop for that line only

		// normal exit
		if i.Err == ErrQuit {
			return
		}
		// all other errors, print and reset error
		if i.Err != nil {
			fmt.Fprintf(i.writer, "%s%s%s\n", ColorRed, i.Err.Error(), ColorOff)
			i.Err = nil
		}

		fmt.Fprint(i.writer, i.Prompt())

		if !linescan.Scan() {
			// End of entry or interrupt
			return
		}

		i.scanner = bufio.NewScanner(strings.NewReader(linescan.Text()))
		i.scanner.Split(i.newSplitFunction())

		i.Run()

	}

}
