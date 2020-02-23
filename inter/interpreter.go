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
	i.rs.push(i.ip + 1) // push return address (up to ';' to pop it)
	i.ip = i.mem[i.ip]  // jump to the dereferenced address
	err = i.interpret() // recurse on the dereferenced cfa
	if err != nil {
		return err
	}

	return err
}
