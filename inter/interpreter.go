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
	// next address to interpret
	ip int
	// here : next free cell in the memory/dictionnary
	here int
	// base used for input/output of numbers
	base int
	// lastNfa, lastPrimitiveNfa
	lastNfa, lastPrimitiveNfa int
}

// NewInterpreter constructor.
func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.writer = os.Stdout

	i.ds, i.rs = newStack(), newStack()

	i.mem = []int{}
	i.words = make(map[int]*word)

	i.alloc(1)
	i.base = 10

	i.initPrimitives()
	return i
}

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
	nfa, err := i.lookup(token)
	if err == nil {
		if i.compileMode {
			return i.compile(nfa + 1)
		}
		i.ip = nfa + 1 // ip points to the cfa of the token.
		return i.interpret()
	}

	// read token as number.
	num, err := strconv.ParseInt(token, i.base, 64)
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
func (i *Interpreter) lookup(token string) (nfa int, err error) {
	return i.lookupFrom(i.lastNfa, token)
}

// lookup only among primitives
func (i *Interpreter) lookupPrimitive(token string) (nfa int, err error) {
	return i.lookupFrom(i.lastPrimitiveNfa, token)
}

// lookup from the lastnfa provided.
func (i *Interpreter) lookupFrom(lastnfa int, token string) (nfa int, err error) {

	// start of search with provied nfa
	nfa = lastnfa
	prevnfa := i.mem[nfa]

	// loop until found or no previous lfa
	for nfa > 0 {
		w := i.words[nfa]
		//fmt.Println("Testing : ", nfa, w)
		if w != nil && w.name == token {
			return nfa, nil
		}

		nfa = prevnfa
		prevnfa = i.mem[nfa]
	}
	return 0, ErrWordNotFound(token)
}

// dump
func (i *Interpreter) dump() {
	fmt.Printf("\n%+v\n", i)
}

// dump
func (i *Interpreter) dumpmem() {
	fmt.Println("Memory dump, size =  ", len(i.mem))
	for k, v := range i.mem {
		fmt.Printf("\t%4d: %8d\n", k, v)
	}
}

// dump
func (i *Interpreter) dumpwords() {
	fmt.Println("Words dumps, size = ", len(i.words))
	for k, w := range i.words {
		fmt.Printf("\t%4d:%+v\n", k, w)
	}
}

// compile the provided cfa on top of the dictionnary
// If it is immediate, call interpret, BUT STAY in compile mode !
func (i *Interpreter) compile(wcfa int) error {

	// check if word is immediate
	w, ok := i.words[wcfa-1]
	if !ok {
		return ErrInvalidCfa(wcfa)
	}
	if w.immediate {
		fmt.Println("Processing an immediate word", w)
		i.ip = wcfa
		return i.interpret()
	}

	i.alloc(1)
	i.mem[i.here-1] = wcfa
	return nil
}

// compile a litteral number
func (i *Interpreter) compileNum(num int) error {
	nfalitt, err := i.lookupPrimitive("LITTERAL")
	if err != nil {
		panic("LITTERAL not defined as primitive ?")
	}
	// write cfa of "litteral" and number
	i.alloc(2)
	i.mem[i.here-2], i.mem[i.here-1] = nfalitt+1, num
	return nil
}

// Interpret the word whose cfa is pointed by the ip pointer
func (i *Interpreter) interpret() (err error) {

	if i.isPrimitive() {
		return i.interpretPrim()
	}

	// compound word
	i.rs.push(i.ip + 1) // push return address (up to ; to pop it)
	i.ip = i.mem[i.ip]  // jump to the dereferenced address
	err = i.interpret() // recurse on the dereferenced cfa
	if err != nil {
		return err
	}
	return err
}
