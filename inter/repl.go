package inter

import "fmt"

// Color definitions.
const (
	ColorRed   = "\033[0;31m"
	ColorGreen = "\033[1;32m"
	ColorBlue  = "\033[1;34m"

	ColorOff = "\033[m"
)

// Prompt the user for entry
func (i *Interpreter) Prompt() string {
	if i.compileMode {
		return fmt.Sprintf("%scompile:%s ", ColorBlue, ColorOff)
	}
	return fmt.Sprintf("%s%d>%s ", ColorGreen, len(i.ds.data), ColorOff)
}

// Repl is the main Read-Evaluate-Print-Loop
func (i *Interpreter) Repl() {

	for {
		fmt.Fprint(i.writer, i.Prompt())
		err := i.Run()
		if err == ErrQuit {
			return
		}
		if err != nil {
			fmt.Fprintf(i.writer, "%s%s%s\n", ColorRed, err.Error(), ColorOff)
		}

	}
}
