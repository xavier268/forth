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
func (i *Interpreter) interpretPrim() (err error) {

	nfa := i.ip - 1
	w, ok := i.words[nfa]
	if !ok {
		return ErrNotPrimitive
	}
	// define all primitive behaviors in interpret mode
	switch w.name {

	case ".":
		n, err := i.ds.pop()
		if err != nil {
			return err
		}
		fmt.Fprintf(i.writer, " %s", strconv.FormatInt(int64(n), i.base))
		i.ip++
		return nil
	case "+":
		n, err := i.ds.pop()
		if err != nil {
			return err
		}
		nn, err := i.ds.pop()
		if err != nil {
			return err
		}
		i.ds.push(n + nn)
		i.ip++
		return nil
	case ":":
		if i.compileMode {
			panic("invalid call to ':' in compile mode")
		}

		fmt.Println("Switching to compile mode")
		i.compileMode = true
		// get next token
		if !i.scanner.Scan() {
			// EOF
			return ErrUnexpectedEndOfLine
		}
		token := i.scanner.Text()

		// create header
		i.createHeader(token)
		return nil

	case ";":
		if i.compileMode {
			// handling the immediate action in compile mode
			fmt.Println("Switching to interpret mode")
			i.alloc(1)
			i.mem[i.here-1] = nfa + 1
			i.compileMode = false
			return nil
		}
		// normal, interpreted mode just pop rs
		i.ip, err = i.rs.pop()
		return err

	default:
		panic("primitive '" + w.name + "' is not implementd")

	}
	// should never get to there ...
}

// createHeader creates a new header in dictionnary.
// updating words, lastNfa and here.
// return nfa of created header.
func (i *Interpreter) createHeader(token string) (nfa int) {
	nfa = i.here
	i.alloc(1)
	i.mem[nfa] = i.lastNfa
	i.words[nfa] = &word{token, false, false}
	i.lastNfa = nfa
	return nfa
}
