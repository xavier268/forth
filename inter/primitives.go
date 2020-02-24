package inter

import (
	"fmt"
	"strconv"
)

func (i *Interpreter) initPrimitives() {

	//           name, immediate
	i.addPrimitive("+", false)
	i.addPrimitive(".", false)
	i.addPrimitive(":", false)
	i.addPrimitive(";", true)
	i.addPrimitive("LITTERAL", false)
	i.addPrimitive("NOOP", false)

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
	// common setting defining the primitive
	nfa := i.ip - 1
	w, ok := i.words[nfa]
	if !ok {
		i.Err = ErrNotPrimitive
		return
	}

	switch w.name {

	case ".":
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		fmt.Fprintf(i.writer, " %s", strconv.FormatInt(int64(n), i.base))

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

	case ":":

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
		// i.alloc(1)
		// i.mem[i.here-1] = cfa

		// switch to compile mode
		fmt.Println("Switching to compile mode")
		i.compileMode = true

	case ";":
		if i.compileMode { // immediate, during compilation

			// write cfa
			i.alloc(1)
			i.mem[i.here-1] = nfa + 1

			// shift back to interpret mode
			fmt.Println("Switching to interpret mode")
			i.compileMode = false
			return
		}
		// normal interpretation in compound word
		// pop one more return address
		ip, err := i.rs.pop()
		if err != nil {
			return // done
		}
		i.ip = ip

	default:
		panic("primitive '" + w.name + "' is not implementd")

	}

	{ // cleanup after normal interpretation of primitive
		var err error
		i.ip, err = i.rs.pop()
		if err != nil {
			return // done
		}
		i.interpret() // loop ...
		return
	}

}
