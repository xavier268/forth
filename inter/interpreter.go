package inter

import (
	"bufio"
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
	//are we currently porcessing a comment ?
	commentMode bool
	// Err contains first interpreter error
	Err error
	// next address to interpret
	ip int

	// lastNfa, lastPrimitiveNfa
	lastNfa, lastPrimitiveNfa int
}

// NewInterpreter constructor.
func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.writer = os.Stdout

	i.ds, i.rs = newStack(), newStack()

	i.words = make(map[int]*word)

	i.initUserVars()
	i.initPrimitives()
	i.initForth()
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

		// remove comments
		if i.commentMode {
			if token == ")" {
				i.commentMode = false
			}
			continue // read next
		}
		if !i.commentMode && token == "(" {
			i.commentMode = true
			continue // read next
		}

		// continue eval ...
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
