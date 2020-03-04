package inter

import (
	"fmt"
	"strconv"
)

// addPrimitive creates and initiates the corresponding primitives.
func (i *Interpreter) addPrimitive(name string) *word {
	// create the words and link them
	w := i.createHeader(name)
	i.lastPrimitiveNfa = i.lastNfa
	w.compil =
		// implements defaults compile behaviour
		func(it *Interpreter) {
			fmt.Println("DEBUG : default compilation of ", it.ip)
			i.mem = append(it.mem, it.ip)
			it.ip = 0
		}
	w.inter =
		// implements default (NOOP) intrepreter behaviour
		func(it *Interpreter) {
			fmt.Println("DEBUG : default (noop) primitive with cfa = ", i.ip)
			it.ip = 0
		}
	return w
}

// define implementation for all primitives.
func (i *Interpreter) initPrimitives() {

	// default finishing function
	// normally, ip=0, and read one more token,
	next := func(i *Interpreter) {
		i.ip = 0
	}

	var w *word

	// bye will terminate the session, exit the repl.
	w = i.addPrimitive("bye")
	w.inter = func(i *Interpreter) {
		i.terminate = true
		i.Err = fmt.Errorf("requested termination")
		next(i)
	}

	// noop will terminate the session, exit the repl.
	w = i.addPrimitive("noop")

	// info will print a dump output
	w = i.addPrimitive("info")
	w.inter = func(i *Interpreter) {
		i.dump()
		next(i)
	}

	// ( a b -- a+b)
	w = i.addPrimitive("+")
	w.inter = func(i *Interpreter) {
		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(a + b)
		next(i)
	}

	// ( a b -- a-b)
	w = i.addPrimitive("-")
	w.inter = func(i *Interpreter) {
		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(b - a)
		next(i)
	}

	// ( a b -- a*b)
	w = i.addPrimitive("*")
	w.inter = func(i *Interpreter) {
		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(a * b)
		next(i)
	}

	// ( n -- ) dot, print ds
	w = i.addPrimitive(".")
	w.inter = func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		fmt.Fprintf(i.writer, " %s", strconv.FormatInt(int64(n), i.getBase()))
		next(i)
	}

	// output following texts until a " word is met,
	// The end of string is marked with a ", even without white spaces.
	// There MUST be a white space after the FIRST "
	w = i.addPrimitive(".\"")
	w.inter = func(i *Interpreter) {
		if i.rs.empty() { // interpreting from repl
			// get the string from the input stream
			token := i.getNextString()
			fmt.Fprintf(i.writer, "%s", token)
		} else { // interpreting from a compound word
			// read the string from memory
			rip, _ := i.rs.pop()                 // return ip
			len := i.mem[rip]                    // get string lenth
			var k int                            // rune pointer
			for k = rip + 1; k <= len+rip; k++ { // retrieve string
				fmt.Fprintf(i.writer, "%s", string(rune(i.mem[k])))
			}
		}
		next(i)
	}
	// in compile mode, will write the cfa, length and string
	// in the dictionnary
	w.compil = func(i *Interpreter) {
		fmt.Printf("DEBUG : Cmode: %v, word: %+v\n", i.compileMode, w)
		token := i.getNextString()
		rtok := []rune(token) // group by rune
		if i.Err != nil {
			return
		}
		fmt.Printf("DEBUG : Cmode: %v, word: %+v\n", i.compileMode, w)
		i.mem = append(i.mem, w.cfa, len(rtok))
		// store the token, rune by rune
		for _, r := range rtok {
			i.mem = append(i.mem, int(r))
		}
		next(i)
	}

	// ( -- addr) addr of the BASE user variable.
	w = i.addPrimitive("base")
	w.inter = func(i *Interpreter) {
		i.Err = i.ds.push(UVBase)
		next(i)
	}

	// enter into intrepretation mode, immediate word
	w = i.addPrimitive("[")
	w.immediate = true
	w.inter = func(i *Interpreter) {
		i.compileMode = false
		next(i)
	}

	// enter into compil mode
	w = i.addPrimitive("]")
	w.inter = func(i *Interpreter) {
		i.compileMode = true
		next(i)
	}

	// make last word immediate
	w = i.addPrimitive("immediate")
	w.inter = func(i *Interpreter) {
		i.words[i.lastNfa].immediate = true
		next(i)
	}

	// ( n -- ) comma, add n in the dictionnary
	w = i.addPrimitive(",")
	w.inter = func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.mem = append(i.mem, n)
		next(i)
	}

	// (n -- ) reserve n cells in the dictionnary
	// enter into intrepretation mode, immediate word
	w = i.addPrimitive("allot")
	w.inter = func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.alloc(n)
		next(i)
	}

	// ( -- addr ) get the address of the first availbale cell of the memory.
	// CAUTION : the memory at 'here' and beyond is NOT ACCESSIBLE unless allocated.
	w = i.addPrimitive("here")
	w.inter = func(i *Interpreter) {
		i.Err = i.ds.push(len(i.mem))
		next(i)
	}

	// (n addr -- ) store n at the given address, assume memory is allocated
	w = i.addPrimitive("!")
	w.inter = func(i *Interpreter) {
		var n, a int
		a, _ = i.ds.pop()
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		if a >= len(i.mem) || a < 0 {
			i.Err = fmt.Errorf("! is trying to store %d at address %d, not accessible", n, a)
			return
		}
		i.mem[a] = n
		next(i)
	}

	// ( addr -- n) get the value n at the address addr
	w = i.addPrimitive("@")
	w.inter = func(i *Interpreter) {
		var a int
		a, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		if a >= len(i.mem) || a < 0 {
			i.Err = fmt.Errorf("@ is trying to access address %d, not accessible", a)
			return
		}
		i.Err = i.ds.push(i.mem[a])
		next(i)
	}

	// ( -- ) emit carriage return
	w = i.addPrimitive("cr")
	w.inter = func(i *Interpreter) {
		fmt.Fprintln(i.writer)
		next(i)
	}

	// ( rune -- ) emit the provided rune
	w = i.addPrimitive("emit")
	w.inter = func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		fmt.Fprintf(i.writer, "%s", string(rune(n)))
		next(i)
	}

	// ( n -- ) drop to of stack
	w = i.addPrimitive("drop")
	w.inter = func(i *Interpreter) {
		_, i.Err = i.ds.pop()
		next(i)
	}

	// ( n -- n n ) dup to of stack
	w = i.addPrimitive("dup")
	w.inter = func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n, n)
		next(i)
	}

	// ( n1 n2 n3 -- n2 n3 n1) rotate the stack
	w = i.addPrimitive("rot")
	w.inter = func(i *Interpreter) {
		var n1, n2, n3 int
		n3, _ = i.ds.pop()
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n2, n3, n1)
		next(i)
	}

	// ( n1 n2  -- n1 n2 n1) over the stack
	w = i.addPrimitive("over")
	w.inter = func(i *Interpreter) {
		var n1, n2 int
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n1, n2, n1)
		next(i)
	}

	// ( n1 n2  -- n2 n1) swap the stack
	w = i.addPrimitive("swap")
	w.inter = func(i *Interpreter) {
		var n1, n2 int
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n2, n1)
		next(i)
	}

	// ( -- r) pop rs into ds
	w = i.addPrimitive("r>")
	w.inter = func(i *Interpreter) {
		var r int
		r, i.Err = i.rs.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(r)
		next(i)
	}

	// ( r -- ) push r into rs
	w = i.addPrimitive(">r")
	w.inter = func(i *Interpreter) {
		var r int
		r, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.rs.push(r)
		next(i)
	}

	// (-- r) push top of rs to ds
	w = i.addPrimitive("r@")
	w.inter = func(i *Interpreter) {
		var r int
		r, i.Err = i.rs.top()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(r)
		next(i)
	}

	// ( -- ) forget <word> : forget the specified word
	// and all the following content, whatever the vocabulary.
	w = i.addPrimitive("forget")
	w.inter = func(i *Interpreter) {

		token := i.getNextString()
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
		next(i)
	}

	// start compiling a compound word
	// do not write the cfa of : in the dictionnary.
	w = i.addPrimitive(":")
	w.inter = func(i *Interpreter) {
		token := i.getNextString()
		if i.Err != nil {
			return
		}
		// create header
		i.createHeader(token)
		// switch to compile mode
		// fmt.Println("Switching to compile mode")
		i.compileMode = true
		next(i)
	}

	// end compiling a compound word, immediate word
	// write the cfa of ; in the dictionnary.
	w = i.addPrimitive(";")
	w.immediate = true
	w.inter = func(i *Interpreter) {
		if i.compileMode { // immediate, during compilation
			// write cfa
			i.mem = append(i.mem, w.cfa)
			// shift back to interpret mode
			// fmt.Println("Switching to interpret mode")
			i.compileMode = false
			i.ip = 0
		}
		// normal interpretation in compound word
		// pop return address
		ip, err := i.rs.pop()
		if err != nil {
			return // done
		}
		i.ip = ip
	}
	/*
		// ( -- n) go get the number that follows and put it on stack
		w = i.addPrimitive("literal")
		w.inter = func(i *Interpreter) {
			if i.rs.empty() {
				i.Err = fmt.Errorf("you cannot use 'literal' in this context")
				return
			}
			// the number is in the next cell
			// get it, and skip it
			i.ds.push(i.mem[i.ip+1])
			i.ip += 2
		}
	*/
}
