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
	// Err contains first interpreter error
	Err error
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
func (i *Interpreter) Run() {
	for {
		if !i.scanner.Scan() {
			// EOF
			return
		}
		token := i.scanner.Text()
		i.Eval(token)
		if i.Err != nil {
			i.Abort()
			return
		}

	}
}

// Eval evaluates token.
func (i *Interpreter) Eval(token string) {

	if i.Err != nil {
		return
	}

	// DEBUG
	fmt.Println("Evaluating : ", token)

	// lookup token in dictionnary
	nfa := i.lookup(token)
	cfa := nfa + 1
	if i.Err == nil {
		if i.compileMode {
			i.compile(cfa)
			return
		}
		i.ip = cfa // ip points to the cfa of the token.
		i.interpret()
		return
	}

	// clear token not found error
	i.Err = nil

	// read token as number.
	num, err := strconv.ParseInt(token, i.base, 64)
	if err != nil {
		i.Err = ErrWordNotFound(token)
		return
	}

	if i.compileMode {
		i.compileNum(int(num))
		return
	}
	i.ds.push(int(num))
	return

}

// compile the provided cfa on top of the dictionnary
// If it is immediate, call interpret, BUT STAY in compile mode !
func (i *Interpreter) compile(wcfa int) {

	if i.Err != nil {
		return
	}

	// check if word is immediate
	w, ok := i.words[wcfa-1]
	if !ok {
		i.Err = ErrInvalidCfa(wcfa)
		return
	}
	if w.immediate {
		fmt.Println("Processing an immediate word", w)
		i.ip = wcfa
		i.interpret()
		return
	}

	i.alloc(1)
	i.mem[i.here-1] = wcfa
	return
}

// compile a litteral number
func (i *Interpreter) compileNum(num int) {

	if i.Err != nil {
		return
	}

	nfalitt := i.lookupPrimitive("LITTERAL")
	if i.Err != nil {
		panic("LITTERAL not defined as primitive ?")
	}
	// write cfa of "litteral" and number
	i.alloc(2)
	i.mem[i.here-2], i.mem[i.here-1] = nfalitt+1, num
	return
}

// Interpret the word whose cfa is pointed by the ip pointer
func (i *Interpreter) interpret() {

	if i.Err != nil {
		return
	}
	// primitive, don't dereference a pseudo cfa
	if i.isPrimitive() {
		i.interpretPrim()
		return
	}

	// compound word
	i.rs.push(i.ip + 1) // push next address on return stack
	i.ip = i.mem[i.ip]  // jump to the dereferenced address
	i.interpret()       // recurse on the dereferenced cfa
	return
}
