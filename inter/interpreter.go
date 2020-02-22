package inter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Interpreter for forth
type Interpreter struct {
	scanner *bufio.Scanner
	writer  io.Writer
	// Data stack, return stack
	ds, rs *stack
	// Mem is the main memory,
	// where the dictionnary lives
	mem []int
	// map nfa to the string details and flags
	words map[int]*word
	// CompileMode (or interpret) mode ?
	compileMode bool
	err         error
}

type word struct {
	name      string
	immediate bool
	smudge    bool
	primitive bool
}

// NewInterpreter constructor.
func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.writer = os.Stdout
	i.scanner = bufio.NewScanner(os.Stdin)
	i.scanner.Split(bufio.ScanWords)
	i.ds, i.rs = newStack(), newStack()

	i.mem = []int{}
	i.words = make(map[int]*word)

	i.initUserVar()
	i.initPrimitives()
	return i
}

// Run the interpreter, until eof
func (i *Interpreter) Run() error {
	for {
		if !i.scanner.Scan() {
			// EOF
			return nil
		}
		token := i.scanner.Text()
		err := i.Eval(token)
		if err != nil {
			i.Abort()
			return err
		}

	}
}

// Abort reset stacks and interpreter
func (i *Interpreter) Abort() {
	i.ds.clear()
	i.rs.clear()
}

// Eval evaluates token.
func (i *Interpreter) Eval(token string) error {

	// lookup token in dictionnary

	// DEBUG
	fmt.Println("Evaluating : ", token)

	panic("not implemented")

}

// Lookup most recent token in disctionnary, using the chain of lfa.
func (i *Interpreter) Lookup(token string) (nfa int, err error) {

	for nfa = i.mem[UVLastNfa]; i.mem[nfa+1] != 0; nfa = i.mem[nfa+1] {
		if w := i.words[nfa]; w != nil && w.name == token {
			return nfa, nil
		}
	}
	return 0, ErrWordNotFound

}
