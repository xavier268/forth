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
	// including implementations details for primitives
	words map[int]*word
	// CompileMode (or interpret) mode ?
	compileMode bool
	// reading string or normal token scan ?
	readingString bool
	// set it to terminate the repl loop
	terminate bool
	// Err contains first interpreter error
	Err error

	// lastNfa, lastPrimitiveNfa
	lastNfa, lastPrimitiveNfa int
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
	//i.initForth()
	return i
}

// Run the interpreter, until eof or another error
func (i *Interpreter) Run() {

	for !i.terminate && i.Err == nil {

		// === read and process next token
		i.ip = 0
		st := i.getNextToken()
		fmt.Printf("DEBUG : just read token : %+v\n", st)
		if i.terminate || i.Err != nil || st.t == errorT {
			if i.Err == nil {
				i.Err = st.err
			}
			return // back to repl or finished, no more token
		}

		// === handle numbers
		if st.t == numberT {
			if i.compileMode { // compile
				cfalit := 1 + i.lookupPrimitive("literal")
				i.mem = append(i.mem, cfalit, st.v)
				i.ip = 0
				continue // getNextToken
			} else { // interpret
				i.ds.push(st.v)
				i.ip = 0
				continue // getNextToken
			}
		}

		// === handle primitives,
		if st.t == primitiveT {
			w := i.words[st.v-1] // read word is indexed on NFA, but st contains CFA !!
			fmt.Printf("DEBUG : about to eval primitive %+v\n", w)
			// set ip to the primitive cfa value
			i.ip = w.cfa
			if i.compileMode && !w.immediate {
				w.compil()
			} else {
				w.inter()
			}
			continue // getNextToken
		}

		// handle compound or primitive - only now is ip significant
		if st.t == compoundT {
			i.ip = st.v // cfa to be executed
			i.eval()
		}
	}
	if i.Err == io.EOF {
		// ignore EOF when returning from Run
		i.Err = nil
	}
}

// navigate the interpreter, using ip and rs
// it is up to the primitives to update ip and rs
// ip is set on entrance
func (i *Interpreter) eval() {

	for i.ip != 0 && i.Err == nil {

		fmt.Printf("DEBUG: evaluating ip : %d -> %d\n", i.ip, i.mem[i.ip])

		// dereference and push rs
		i.rs.push(i.ip + 1)
		i.ip = i.mem[i.ip]

		// handle primitives, they need to manage rs and ip
		// default is to increment ip and to touch rs
		if i.ip <= 1+i.lastPrimitiveNfa {
			w := i.words[i.ip-1]
			if i.compileMode && !w.immediate {
				w.compil()
			} else {
				w.inter()
			}
		}
	}
}

/*
// Eval evaluates token.
func (i *Interpreter) Eval(token string) {

	if i.Err != nil {
		return
	}

	// DEBUG
	// fmt.Println("Evaluating : ", token)

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
	num, err := strconv.ParseInt(token, i.getBase(), 64)
	if err != nil {
		i.Err = ErrWordNotFound(token)
		return
	}
	// compile numbre ...
	if i.compileMode {
		i.compileNum(int(num))
		return
	}
	// ... or push it to stack if interpret mode
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
	if w == nil || !ok {
		i.Err = ErrInvalidCfa(wcfa)
		return
	}
	if w.immediate {
		// fmt.Println("Processing an immediate word : '", w.name, "'")
		i.ip = wcfa
		i.interpret()
		return
	}

	// handle special compile mode behaviours,
	// or use defaut behaviour

	i.compilePrim(wcfa, w)
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

*/
