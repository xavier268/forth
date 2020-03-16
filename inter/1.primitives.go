package inter

import (
	"fmt"
	"strconv"
)

// addPrimitive creates and initiates a word
// for the corresponding primitives.
func (i *Interpreter) addPrimitive(name string) (pcode int) {

	// create the words and link them
	w := i.createHeader(name)
	// create a pcode from the nfa,
	// compile the pseudo code in the cfa
	pcode = -w.nfa
	i.mem = append(i.mem, pcode)
	return pcode
}

// add an immediate primitive
func (i *Interpreter) addPrimitiveImmediate(name string) (pcode int) {
	// create the words and link them
	w := i.createHeader(name)
	// make it immediate
	w.immediate = true
	// create a pcode from the nfa,
	// compile the pseudo code in the cfa
	pcode = -w.nfa
	i.mem = append(i.mem, pcode)
	return pcode
}

// default move of the ip pointer at the end of a primitive
func (i *Interpreter) moveIP() {
	// if interpret and non empty rs, increment ip
	if i.ip != 0 && !i.compileMode && !i.rs.empty() {
		i.ip++
		return
	}
	i.ip = 0
}

// define implementation for all primitives.
func (i *Interpreter) initPrimitives() {

	var pcfa int
	i.code = NewPrimCode(
		func(ii *Interpreter) {
			fmt.Println("DEBUG : Calling default interpret primitive")
			i.moveIP()
		},
		func(i2 *Interpreter) {
			fmt.Println("DEBUG : Calling default compile primitive")
			i.moveIP()
		})

	// bye will terminate the session, exit the repl.
	pcfa = i.addPrimitive("bye")
	i.code.addInter(pcfa, func(i *Interpreter) {
		i.terminate = true
		i.Err = fmt.Errorf("requested termination")
		i.moveIP()
	})

	// noop will terminate the session, exit the repl.
	i.addPrimitive("noop")

	// info will print a dump output
	pcfa = i.addPrimitive("info")
	i.code.addInter(pcfa, func(i *Interpreter) {

		i.dump()
		i.moveIP()
	})

	// ( a b -- a+b)
	pcfa = i.addPrimitive("+")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(a + b)
		i.moveIP()
	})

	// ( a b -- a-b)
	pcfa = i.addPrimitive("-")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(b - a)
		i.moveIP()
	})

	// ( a b -- a*b)
	pcfa = i.addPrimitive("*")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(a * b)
		i.moveIP()
	})

	// ( n -- ) dot, print ds
	pcfa = i.addPrimitive(".")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		fmt.Fprintf(i.writer, " %s", strconv.FormatInt(int64(n), i.getBase()))
		i.moveIP()
	})

	// output following texts until a " word is met,
	// The end of string is marked with a ", even without white spaces.
	// There MUST be a white space after the FIRST "
	pcfa = i.addPrimitive(".\"")
	i.code.addInter(pcfa, func(i *Interpreter) {

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
		i.moveIP()
	})
	// in compile mode, will write the cfa, length and string
	// in the dictionnary
	i.code.addCompil(pcfa, func(i *Interpreter) {

		fmt.Printf("DEBUG : Cmode: %v, word: %+v\n", i.compileMode, pcfa)
		token := i.getNextString()
		rtok := []rune(token) // group by rune
		if i.Err != nil {
			return
		}
		fmt.Printf("DEBUG : Cmode: %v, word: %+v\n", i.compileMode, pcfa)
		i.mem = append(i.mem, pcfa, len(rtok))
		// store the token, rune by rune
		for _, r := range rtok {
			i.mem = append(i.mem, int(r))
		}
		i.moveIP()
	})

	// ( -- addr) addr of the BASE user variable.
	pcfa = i.addPrimitive("base")
	i.code.addInter(pcfa, func(i *Interpreter) {

		i.Err = i.ds.push(UVBase)
		i.moveIP()
	})

	// enter into intrepretation mode, immediate word
	pcfa = i.addPrimitiveImmediate("[")
	i.code.addInter(pcfa, func(i *Interpreter) {

		i.compileMode = false
		i.moveIP()
	})

	// enter into compil mode
	pcfa = i.addPrimitive("]")
	i.code.addInter(pcfa, func(i *Interpreter) {

		i.compileMode = true
		i.moveIP()
	})

	// make last word immediate
	pcfa = i.addPrimitive("immediate")
	i.code.addInter(pcfa, func(i *Interpreter) {

		i.words[i.lastNfa].immediate = true
		i.moveIP()
	})

	// ( n -- ) comma, add n in the dictionnary
	pcfa = i.addPrimitive(",")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.mem = append(i.mem, n)
		i.moveIP()
	})

	// (n -- ) reserve n cells in the dictionnary
	// enter into intrepretation mode, immediate word
	pcfa = i.addPrimitive("allot")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.alloc(n)
		i.moveIP()
	})

	// ( -- addr ) get the address of the first availbale cell of the memory.
	// CAUTION : the memory at 'here' and beyond is NOT ACCESSIBLE unless allocated.
	pcfa = i.addPrimitive("here")
	i.code.addInter(pcfa, func(i *Interpreter) {

		i.Err = i.ds.push(len(i.mem))
		i.moveIP()
	})

	// (n addr -- ) store n at the given address, assume memory is allocated
	pcfa = i.addPrimitive("!")
	i.code.addInter(pcfa, func(i *Interpreter) {

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
		i.moveIP()
	})

	// ( addr -- n) get the value n at the address addr
	pcfa = i.addPrimitive("@")
	i.code.addInter(pcfa, func(i *Interpreter) {

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
		i.moveIP()
	})

	// ( -- ) emit carriage return
	pcfa = i.addPrimitive("cr")
	i.code.addInter(pcfa, func(i *Interpreter) {

		fmt.Fprintln(i.writer)
		i.moveIP()
	})

	// ( rune -- ) emit the provided rune
	pcfa = i.addPrimitive("emit")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		fmt.Fprintf(i.writer, "%s", string(rune(n)))
		i.moveIP()
	})

	// ( n -- ) drop to of stack
	pcfa = i.addPrimitive("drop")
	i.code.addInter(pcfa, func(i *Interpreter) {

		_, i.Err = i.ds.pop()
		i.moveIP()
	})

	// ( n -- n n ) dup to of stack
	pcfa = i.addPrimitive("dup")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n, n)
		i.moveIP()
	})

	// ( n1 n2 n3 -- n2 n3 n1) rotate the stack
	pcfa = i.addPrimitive("rot")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n1, n2, n3 int
		n3, _ = i.ds.pop()
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n2, n3, n1)
		i.moveIP()
	})

	// ( n1 n2  -- n2 n1) swap the stack
	pcfa = i.addPrimitive("swap")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n1, n2 int
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n2, n1)
		i.moveIP()
	})

	// ( -- r) pop rs into ds
	pcfa = i.addPrimitive("r>")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var r int
		r, i.Err = i.rs.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(r)
		i.moveIP()
	})

	// ( r -- ) push r into rs
	pcfa = i.addPrimitive(">r")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var r int
		r, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.rs.push(r)
		i.moveIP()
	})

	// (-- r) push top of rs to ds
	pcfa = i.addPrimitive("r@")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var r int
		r, i.Err = i.rs.top()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(r)
		i.moveIP()
	})

	// end compiling a compound word, immediate word
	// write the pcfa of ; in the dictionnary.
	pcfa = i.addPrimitiveImmediate(";")
	i.code.addInter(pcfa, func(i *Interpreter) {

		if i.compileMode { // immediate, during compilation
			// write cfa
			fmt.Printf("DEBUG : compiling cfa of ; as %d - %+v\n", pcfa, pcfa)
			i.mem = append(i.mem, pcfa)
			// shift back to interpret mode
			fmt.Println("DEBUG : Switching to interpret mode")
			i.compileMode = false
			i.ip = 0
			return // done !
		}
		// normal interpretation in compound word
		// pop return address
		if i.rs.empty() {
			i.ip = 0
			return
		}
		i.ip, i.Err = i.rs.pop()
		if i.Err != nil {
			i.ip = 0
			return // done
		}
	})

	// start compiling a compound word
	// do not write the cfa of : in the dictionnary.
	pcfa = i.addPrimitive(":")
	i.code.addInter(pcfa, func(i *Interpreter) {

		token := i.getNextString()
		if i.Err != nil {
			return
		}
		// create header
		i.createHeader(token)
		// switch to compile mode
		// fmt.Println("Switching to compile mode")
		i.compileMode = true
		i.moveIP()
	})

	// ( -- ) forget <word> : forget the specified word
	// and all the following content, whatever the vocabulary.
	pcfa = i.addPrimitive("forget")
	i.code.addInter(pcfa, func(i *Interpreter) {

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
		i.moveIP()
	})

	// ( n1 n2  -- n1 n2 n1) over the stack
	pcfa = i.addPrimitive("over")
	i.code.addInter(pcfa, func(i *Interpreter) {

		var n1, n2 int
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n1, n2, n1)
		i.moveIP()
	})

	/*
		// ( -- n) go get the number that follows and put it on stack
		w = i.addPrimitive("literal")
		w.inter = func(){ i:= i
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
