package inter

import (
	"errors"
	"fmt"
	"strconv"
)

// interpretPrim based on the cfa pointed to by IP
// It is the core of the interpreter, that implements
// all the primitive interpretation behavior,
// normally in interpret mode, but also in
// compile mode for immediate words.
func (i *Interpreter) interpretPrim() {

	if i.Err != nil {
		return
	}
	// common setting defining the primitive
	nfa := i.ip - 1
	w, ok := i.words[nfa]
	if !ok {
		i.Err = ErrNotPrimitive(nfa)
		return
	}

	switch w.name {

	case "BYE": // Exit program.
		i.Err = ErrQuit
		return

	case "ABORT", "RESET": // Reset REPL, clean stacks.
		i.Abort()
		return

	case "INFO": // dump debugging info
		i.dump()

	case "BASE": // ( -- addr )
		i.ds.push(UVBase)

	case "[": // immediate word
		if !i.compileMode {
			i.Err = ErrWrongContextWord(w.name)
		}
		i.compileMode = false

	case "]":
		if i.compileMode {
			i.Err = ErrWrongContextWord(w.name)
		}
		i.compileMode = true

	case "IMMEDIATE": // lastNFA word will be made immediate
		i.words[i.lastNfa].immediate = true

	case ",": // (n -- ) Add n to the next dictionnary cell, allocating ONE cell.
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		i.mem = append(i.mem, n)

	case "ALLOT": // (n --) Add n cells to the dictionnay.

		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}

		i.alloc(n)

	case "HERE": // ( -- addr ) get the address of the first availbale cell of the memory.
		// CAUTION : the memory at HERE and beyond is NOT ACCESSIBLE unless allocated.

		i.ds.push(len(i.mem))

	case "!": // (n addr --) store n at the given address, if it is allocated
		a, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}

		if a >= len(i.mem) || a < 0 {
			i.Err = ErrInvalidAddr(a)
			return
		}

		i.mem[a] = n

	case "@": // (addr -- n) fetch memory content
		a, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		if a >= len(i.mem) || a < 0 {
			i.Err = ErrInvalidAddr(a)
			return
		}

		i.ds.push(i.mem[a])

	case ".\"":
		// output following texts until a " word is met,
		// The end of string is marked with a ", even without white spaces.
		// There MUST be a white space after the first "

		if i.rs.empty() { // interpreting from repl
			// get the string from the input stream
			token := i.scanNextToken()
			fmt.Fprintf(i.writer, "%s", token)
		} else { // interpreting from a compound word
			// read the string from memory
			rip, _ := i.rs.pop()                 // return ip
			len := i.mem[rip]                    // get string lenth
			var k int                            // rune pointer
			for k = rip + 1; k <= len+rip; k++ { // retrieve string
				fmt.Fprintf(i.writer, "%s", string(rune(i.mem[k])))
			}
			// update new return address
			i.rs.push(k)
		}

	case "CR": // emit carriage return
		fmt.Fprintln(i.writer)

	case "EMIT": // ( char -- ) emit the char
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}

		fmt.Fprintf(i.writer, "%s", string(rune(n)))

	case ".":
		n, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		//fmt.Println("DEBUG : BASE = ", i.getBase())
		fmt.Fprintf(i.writer, " %s",
			strconv.FormatInt(int64(n), i.getBase()))

	case "DROP": // ( n -- n n )
		_, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}

	case "DUP": // ( n -- n n )
		n, err := i.ds.top()
		if err != nil {
			i.Err = err
			return
		}
		i.ds.push(n)
	case "ROT":
		l := len(i.ds.data)
		if l < 3 {
			i.Err = ErrStackUnderflow
			return
		}
		i.ds.data[l-1], i.ds.data[l-2], i.ds.data[l-3] =
			i.ds.data[l-3], i.ds.data[l-1], i.ds.data[l-2]

	case "OVER":
		l := len(i.ds.data)
		if l < 2 {
			i.Err = ErrStackUnderflow
			return
		}
		i.ds.push(i.ds.data[l-2])

	case "SWAP": // (a b -- b a )
		l := len(i.ds.data)
		if l < 2 {
			i.Err = ErrStackUnderflow
			return
		}
		i.ds.data[l-2], i.ds.data[l-1] = i.ds.data[l-1], i.ds.data[l-2]

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

	case "*":
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

		i.ds.push(n * nn)

	case "-": // ( n1 n2 -- "n2-n1") Substract
		n2, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}
		n1, err := i.ds.pop()
		if err != nil {
			i.Err = err
			return
		}

		i.ds.push(n1 - n2)

	case "FORGET": // FORGET <xxx> will remove the xxx word
		// all subsequent dictionnary cells will become unavialable

		token := i.scanNextToken()
		if i.Err != nil {
			return
		}

		nfa2forget := i.lookup(token)
		i.lastNfa = i.mem[nfa2forget]

		if i.Err != nil {
			return
		}
		// Cleanup mem
		i.mem = i.mem[:nfa2forget]
		// Cleanup words that are not accessible anymore
		for nfa2 := range i.words {
			if nfa2 >= nfa2forget {
				// surprinsingly,
				// it is safe to delete in a range,
				// according to go spec !
				delete(i.words, nfa2)
			}
		}

	case ":":

		token := i.scanNextToken()
		if i.Err != nil {
			return
		}

		// create header
		i.createHeader(token)

		// switch to compile mode
		// fmt.Println("Switching to compile mode")
		i.compileMode = true

	case ";":
		if i.compileMode { // immediate, during compilation

			// write cfa
			i.alloc(1)
			i.mem[len(i.mem)-1] = nfa + 1

			// shift back to interpret mode
			// fmt.Println("Switching to interpret mode")
			i.compileMode = false

		}
		// normal interpretation in compound word
		// pop one more return address
		ip, err := i.rs.pop()
		if err != nil {
			return // done
		}
		i.ip = ip

	case "NOOP":
		// do nothing

	case "LITERAL":

		if i.rs.empty() {
			i.Err = ErrWrongContextWord(w.name)
			return
		}

		// the number is pointed by the return stack
		// get it, and points to the following address
		nextip, _ := i.rs.pop()
		i.rs.push(nextip + 1)
		i.ds.push(i.mem[nextip])

	case "<BUILDS":
		token := i.scanNextToken()
		if i.Err != nil {
			return
		}
		// create header
		i.createHeader(token)

		nfad := i.lookupPrimitive("$$DOES$$")

		// store $$DOES$$ cfa, and
		// an empty slot,taht DOES> will fill
		// write -1 in it to generate errors if not filled correctly !
		i.mem = append(i.mem, nfad+1, -1)

		// save the address of the reserved slot
		// in the return stack, just under the top
		next, _ := i.rs.pop()
		i.rs.push(len(i.mem) - 1)
		i.rs.push(next)

		// continue interpreting next <BUILDS instructions

	case "DOES>":
		// expect at least 2 values in return stack.
		if len(i.rs.data) < 2 {
			i.Err = errors.New("you need to use <BUILDS before calling DOES>")
		}
		// get the address of the process that $$DOES$$ will execute
		// store the instructions set2 in the reserved address
		addrProc2, _ := i.rs.pop()
		reservedAddr, _ := i.rs.pop()
		i.mem[reservedAddr] = addrProc2

		// do NOT interpret further at this stage, so force a further pop
		// and send to execution
		i.ip, _ = i.rs.pop()

	case "$$DOES$$":

	case "CONSTANT":

		token := i.scanNextToken()
		if i.Err != nil {
			return
		}

		// create header
		i.createHeader(token)

		// Get the number,
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}

		// compile the number with the $$CONSTANT$$ cfa
		nfa := i.lookupPrimitive("$$CONSTANT$$")
		i.mem = append(i.mem, nfa+1, n)

	case "$$CONSTANT$$":

		if i.rs.empty() {
			i.Err = ErrWrongContextWord(w.name)
			return
		}

		// the number is pointed by the return stack
		// get it, and drop a return stack level
		nextip, _ := i.rs.pop()
		i.ds.push(i.mem[nextip])

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
