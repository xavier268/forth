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
	// next address to interpret
	ip int
	// Mem is the main memory,
	// where the dictionnary lives
	mem []int
	// map NFA to the string details and flags,
	words map[int]*word
	// CompileMode (or interpret) mode ?
	compileMode bool
	// reading string or normal token scan ?
	readingString bool
	// set it to terminate the repl loop
	terminate bool
	// Err contains first interpreter error
	Err error

	// code for primitives is stored here.
	code *PrimCode

	// lastNfa
	lastNfa          int
	lastPrimitiveNfa int
}

// NewInterpreter constructor.
func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.SetWriter(os.Stdout)
	i.SetReader(os.Stdin)

	i.ds, i.rs = newStack(), newStack()
	i.ds.errUnder = fmt.Errorf("data stack underflow")
	i.rs.errUnder = fmt.Errorf("return stack underflow")

	i.words = make(map[int]*word)

	i.initUserVars()
	i.initPrimitives()
	i.initForth()
	return i
}

// Run the interpreter, until eof or another error
func (i *Interpreter) Run() {

	for !i.terminate && i.Err == nil {

		// === read and process next token
		i.ip = 0
		st := i.getNextToken()
		// fmt.Printf("DEBUG : just read token : %+v\n", st)
		if st.t == errorT {
			i.Err = st.err
			continue // back to repl or finished, no more token
		}
		if i.terminate || i.Err != nil {
			continue // // back to repl or finished, no more token
		}

		// === handle numbers
		if st.t == numberT {
			if i.compileMode { // compile number
				cfalit := -i.lookupFrom(i.lastPrimitiveNfa, "literal")
				i.mem = append(i.mem, cfalit, st.v)
				continue // getNextToken
			} else { // interpret number
				i.ds.push(st.v)
				continue // getNextToken
			}
		}

		// handle compound - only now is ip significant
		if st.t == compoundT {
			switch {
			case !i.compileMode || // normal interpretation
				i.words[st.v-1].immediate: // or immediate
				i.rs.push(0) // push repl level
				i.ip = st.v
				i.eval()
				i.ip, _ = i.rs.pop()
			case i.compileMode:
				// normal compilation
				i.mem = append(i.mem, st.v)
			}
		}
	}
	// ignore EOF when returning from Run
	if i.Err == io.EOF {
		i.Err = nil
	}
}

// navigate the interpreter, using ip and rs
// it is up to the primitives to update ip and rs
// ip is set on entrance. Is normally never called with ip = 0.
// mode can be compile & immediate, or interpret.
func (i *Interpreter) eval() {

	for i.ip != 0 && i.Err == nil {

		// fmt.Printf("DEBUG : evaluating ip : %d -> %d, rs: %+v\n", i.ip, i.mem[i.ip], i.rs.data)

		// if pointing to pseudo code, we have a primitive !
		if i.mem[i.ip] < 0 {
			// execute primitive code !
			i.code.do(i, i.mem[i.ip])
		} else {
			//  compound,
			// dereference and push rs
			i.rs.push(i.ip + 1)
			i.ip = i.mem[i.ip]
		}
	}
}
