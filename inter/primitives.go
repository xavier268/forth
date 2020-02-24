package inter

import (
	"fmt"
	"strconv"
)

func (i *Interpreter) initPrimitives() {

	//           name, immediate
	i.addPrimitive("+", false)
	i.addPrimitive(".", false)
	i.addPrimitive(":", true)
	i.addPrimitive(";", true)
	i.addPrimitive("LITTERAL", false)

	// flag last primitive nfa
	i.lastPrimitiveNfa = i.lastNfa
}

func (i *Interpreter) addPrimitive(name string, immediate bool) {

	nfa := i.createHeader(name)
	if immediate {
		i.words[nfa].immediate = true
	}
}

// detect if ip is pointing to a primitive cfa
func (i *Interpreter) isPrimitive() bool {
	return i.ip-1 <= i.lastPrimitiveNfa
}

// interpretPrim based on the cfa pointed to by IP
func (i *Interpreter) interpretPrim() {

	if i.Err != nil {
		return
	}

	nfa := i.ip - 1
	w, ok := i.words[nfa]
	if !ok {
		i.Err = ErrNotPrimitive
		return
	}

	// define all primitive behaviors in interpret mode
	switch w.name {

	case ".":
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		fmt.Fprintf(i.writer, " %s", strconv.FormatInt(int64(n), i.base))
		i.ip, err = i.rs.pop()
		if err != nil {
			panic(err)
		}
		i.interpret()
		return
	case "+":
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		nn, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		i.ds.push(n + nn)
		i.ip, err = i.rs.pop()
		if err != nil {
			panic(err)
		}
		i.interpret()
		return
	case ":":
		if i.compileMode { // immediate in compile mode
			i.ip++
			i.interpret()
			return
		}

		fmt.Println("Switching to compile mode")
		i.compileMode = true
		// get next token
		if !i.scanner.Scan() {
			// EOF
			i.Err = ErrUnexpectedEndOfLine
			return
		}
		token := i.scanner.Text()

		// create header
		i.createHeader(token)

		// add cfa of the definition word, :
		i.alloc(1)
		i.mem[i.here-1] = nfa + 1
		return

	case ";":
		if i.compileMode {
			// handling the immediate action in compile mode
			fmt.Println("Switching to interpret mode")
			i.alloc(1)
			i.mem[i.here-1] = nfa + 1
			i.compileMode = false
			return
		}
		// normal, interpreted mode just pop rs
		i.ip, i.Err = i.rs.pop()
		return

	default:
		panic("primitive '" + w.name + "' is not implementd")

	}
	// should never get to there ...
}
