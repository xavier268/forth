package inter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
	// IP ?
}

// Eval evaluates token.
func (i *Interpreter) Eval(token string) error {

	// DEBUG
	fmt.Println("Evaluating : ", token)

	// lookup token in dictionnary
	_, w, err := i.lookup(token)
	if err == nil {
		if i.compileMode {
			return i.compile(w)
		}
		return i.interpret(w)

	}

	// read token as number.
	num, err := strconv.ParseInt(token, i.mem[UVBase], 64)
	if err != nil {
		return ErrWordNotFound(token)
	}

	if i.compileMode {
		i.compileNum(int(num))
	} else {
		i.ds.push(int(num))
	}
	return nil
}

// lookup most recent token in dictionnary, using the chain of lfa.
func (i *Interpreter) lookup(token string) (nfa int, w *word, err error) {

	// start of search
	nfa = i.mem[UVLastNfa]
	prevnfa := i.mem[nfa]

	// loop until found or no previous lfa
	for nfa > 0 {
		w := i.words[nfa]
		//fmt.Println("Testing : ", nfa, w)
		if w != nil && w.name == token {
			return nfa, w, nil
		}

		nfa = prevnfa
		prevnfa = i.mem[nfa]
	}
	return 0, nil, ErrWordNotFound(token)
}

func (i *Interpreter) dumpmem() {
	fmt.Println("Memory dump, size =  ", len(i.mem))
	for k, v := range i.mem {
		fmt.Printf("\t%4d: %8d\n", k, v)
	}
}

func (i *Interpreter) dumpwords() {
	fmt.Println("Words dumps, size = ", len(i.words))
	for k, w := range i.words {
		fmt.Printf("\t%4d:%+v\n", k, w)
	}
}

func (i *Interpreter) compile(*word) error   { return nil }
func (i *Interpreter) compileNum(int) error  { return nil }
func (i *Interpreter) interpret(*word) error { return nil }
