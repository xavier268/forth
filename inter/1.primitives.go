package inter

import (
	"errors"
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
	colon := -i.lookupFrom(i.lastPrimitiveNfa, ";")
	i.mem = append(i.mem, pcode, colon)
	i.lastPrimitiveNfa = i.lastNfa
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
	colon := -i.lookupFrom(i.lastPrimitiveNfa, ";")
	i.mem = append(i.mem, pcode, colon)
	i.lastPrimitiveNfa = i.lastNfa
	return pcode
}

// default move of the ip pointer at the end of a primitive
// primitive code is NEVER called directly, always as part of eval.
func (i *Interpreter) moveIP() {
	if i.ip != 0 {
		i.ip++
	}
}

// define implementation for all primitives.
func (i *Interpreter) initPrimitives() {

	// initialize an empty PrimCode structure.
	i.code = NewPrimCode(
		func(i1 *Interpreter) {
			fmt.Printf("WARNING : Calling default interpret primitive, ip:%d->%d\n",
				i1.ip, i1.mem[i1.ip])
			i.moveIP()
		},
		func(i2 *Interpreter) {
			fmt.Printf("WARNING : Calling default compile(immediate) primitive, ip:%d->%d\n",
				i2.ip, i2.mem[i2.ip])
			i.moveIP()
		})

	// end compiling a compound word, immediate word
	// write the pcfa of ; in the dictionnary.
	// needs to be defined early, it is used even for primitives.
	{
		pcfa := i.addPrimitiveImmediate(";")
		i.code.addInter(pcfa, func(i *Interpreter) {
			// normal interpretation in compound word
			// pop return address, leaving 0 if stack is empty.
			i.ip, i.Err = i.rs.pop()
			// reset on error OR if rs is empty
			if i.Err != nil {
				fmt.Printf("WARNING : resetting error ? : %v, (ip:%d, rs:%+v)\n", i.Err, i.ip, i.rs.data)
				i.Err = nil
				i.ip = 0
			}
		})
		i.code.addCompil(pcfa, func(i *Interpreter) {
			// immediate, during compilation
			// write cfa
			// fmt.Printf("DEBUG : compiling cfa of ; as %d \n", pcfa)
			i.mem = append(i.mem, pcfa)
			// shift back to interpret mode
			// fmt.Println("DEBUG : Switching to interpret mode")
			i.compileMode = false
			i.ip = 0 // to ask for a new token ...
			return   // done !
		})
	}

	// ( -- n) go get the number that follows and put it on stack
	i.code.addInter(i.addPrimitive("literal"),
		func(i *Interpreter) {
			// the number is in the ip+1 cell
			// get it, and skip it
			i.ip++
			i.ds.push(i.mem[i.ip])
			i.moveIP()
		})

	// ( -- ) $$does$$ is an internal word, that will be called
	// from the object to run the executable part defined with does>
	// the executable code will be passed the pfa of the built object
	// on the data stack
	i.code.addInter(i.addPrimitive("$$does$$"),
		func(i *Interpreter) {
			var cadr, pfa int
			if i.rs.empty() {
				i.Err = fmt.Errorf("the return stack is inconsistent(empty) while interpreting $$does$$")
			}
			fmt.Printf("DEBUG : executig $$does$$. ip:%d, rs:%+v\n", i.ip, i.rs.data)

			// get executable code adress
			cadr = 1 + i.ip

			// push pfa on ds
			pfa = 1 + cadr
			i.ds.push(pfa)

			// jump to executable code
			i.ip = i.mem[cadr]

			// return will be handled by executable code ...
		})

	// ( -- ) does> writes the address of the executable code
	// into the reserved sapce prepared by <builds
	// whose address is in the return stack.
	i.code.addInter(i.addPrimitive("does>"),
		func(i *Interpreter) {
			fmt.Printf("DEBUG : entering 'does>' with ip:%d,rs%+v\n", i.ip, i.rs.data)

			if i.rs.empty() {
				i.Err = errors.New("cannot use 'does>' in this context")
				return
			}
			// get the address of the executable code
			addrProc2, _ := i.rs.pop()

			// get the reserverd word to store that address
			addrResevd, _ := i.rs.pop()
			// store the instructions set of Proc2
			// (immediately following current ip) in the reserved address
			fmt.Printf("DEBUG : writing %d in addr %d\n", addrProc2, addrResevd)
			i.mem[addrResevd] = addrProc2

			// return, do not execute further !
			i.ip, i.Err = i.rs.pop()

		})

	// ( -- ) <builds creates a header, and add after it,
	// the pcode for $$does$$ and the address of the executable code
	i.code.addInter(i.addPrimitive("<builds"),
		func(i *Interpreter) {
			token := i.getNextString()
			if i.Err != nil {
				i.Err = errors.New("expected token name but got EOF")
				return
			}
			if len(i.rs.data) < 2 {
				i.Err = errors.New("'<builds' cannot be called in this context")
				return
			}
			// create header
			i.createHeader(token)

			// get the pseudo-cfa of $$does$$
			cfad := -i.lookupFrom(i.lastPrimitiveNfa, "$$does$$")

			// store $$DOES$$ pcfa, and
			// an empty slot,that DOES> will fill
			// write -1000 in it to generate errors if not filled correctly !
			i.mem = append(i.mem, cfad, -1000)

			// save the address of the reserved slot
			// in the return stack, just under the top
			top, _ := i.rs.pop()
			i.rs.push(len(i.mem) - 1)
			i.rs.push(top)

			// continue interpreting next <BUILDS instructions
			i.moveIP()
		})

	// ( -- ) branch execution to an offset address,
	// offset is read from the next cell
	// use offset = 0 for normal execution flow.
	i.code.addInter(i.addPrimitive("branch"),
		func(i *Interpreter) {
			//fmt.Printf("DEBUG : entering branch @ ip : %d, rs: %+v\n", i.ip, i.rs.data)

			var next, offset int
			if len(i.rs.data) < 2 || i.ip == 0 {
				i.Err = errors.New("you cannot use 'branch' in this context")
				return
			}
			next, _ = i.rs.pop()
			offset = i.mem[next]
			//fmt.Println("DEBUG : offset is ", offset)
			if i.Err != nil {
				return
			}
			i.rs.push(offset + next + 1) // skip radr cell for normal flow
			//fmt.Printf("DEBUG : exiting branch @ ip : %d, rs: %+v\n", i.ip, i.rs.data)
			i.moveIP()
		})

	// ( flag -- ) 0branch execution to an offset address if flag is 0
	// offset is read from the next cell
	// use offset = 0 for normal execution flow.
	i.code.addInter(i.addPrimitive("0branch"),
		func(i *Interpreter) {
			//fmt.Printf("DEBUG : entering branch @ ip : %d, rs: %+v\n", i.ip, i.rs.data)

			var next, offset, flag int
			if len(i.rs.data) < 2 || i.ip == 0 {
				i.Err = errors.New("you cannot use 'branch' in this context")
				return
			}
			flag, i.Err = i.ds.pop()
			next, _ = i.rs.pop()
			if flag == 0 {
				offset = i.mem[next]
			} else {
				offset = 0
			}
			//fmt.Println("DEBUG : offset is ", offset)
			if i.Err != nil {
				return
			}
			i.rs.push(offset + next + 1) // skip radr cell for normal flow
			//fmt.Printf("DEBUG : exiting branch @ ip : %d, rs: %+v\n", i.ip, i.rs.data)
			i.moveIP()
		})

	// bye will terminate the session, exit the repl.
	i.code.addInter(i.addPrimitive("bye"), func(i *Interpreter) {
		i.terminate = true
		i.Err = fmt.Errorf("requested termination")
		i.moveIP()
	})

	// ( -- cfa ) tick puts on the data stack the CFA of the following word
	i.code.addInter(i.addPrimitive("'"), func(i *Interpreter) {
		token := i.getNextString()
		if i.Err != nil {
			return
		}
		cfa := 1 + i.lookup(token)
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(cfa)
		i.moveIP()
	})

	// start compiling a compound word
	// do not write the cfa of : in the dictionnary.
	i.code.addInter(i.addPrimitive(":"), func(i *Interpreter) {
		token := i.getNextString()
		if i.Err != nil {
			return
		}
		// create header
		i.createHeader(token)
		// switch to compile mode
		// fmt.Println("DEBUG : Switching to compile mode")
		i.compileMode = true
		i.ip = 0
		i.moveIP()
	})

	// noop does nothing.
	i.code.addInter(i.addPrimitive("noop"), func(i *Interpreter) {
		i.moveIP()
	})

	// ( cfa -- ) execute the cfa that is on the data stack.
	i.code.addInter(i.addPrimitive("execute"), func(i *Interpreter) {
		i.ip, i.Err = i.ds.pop()
		i.eval()
		i.moveIP()
	})

	// ( nfa -- ) print the name corresponding to the nfa
	// no spacing, no cr included.
	i.code.addInter(i.addPrimitive(".name"), func(i *Interpreter) {
		var nfa int
		nfa, i.Err = i.ds.pop()
		w, ok := i.words[nfa]
		if ok && w != nil && w.name != "" {
			fmt.Fprint(i.writer, w.name)
		}
		i.moveIP()
	})

	// info will print a dump output
	i.code.addInter(i.addPrimitive("info"), func(i *Interpreter) {
		i.dump()
		i.moveIP()
	})

	// ( a b -- a+b)
	i.code.addInter(i.addPrimitive("+"), func(i *Interpreter) {
		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(a + b)
		i.moveIP()
	})

	// ( a b -- a-b)
	i.code.addInter(i.addPrimitive("-"), func(i *Interpreter) {
		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(b - a)
		i.moveIP()
	})

	// ( a b -- a*b)
	i.code.addInter(i.addPrimitive("*"), func(i *Interpreter) {
		var a, b int
		a, _ = i.ds.pop()
		b, i.Err = i.ds.pop()
		i.ds.push(a * b)
		i.moveIP()
	})

	// ( n -- ) dot, print ds
	// no space added
	i.code.addInter(i.addPrimitive("."), func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		fmt.Fprintf(i.writer, "%s", strconv.FormatInt(int64(n), i.getBase()))
		i.moveIP()
	})

	{
		// output following texts until a " word is met,
		// The end of string is marked with a ", even without white spaces.
		// There MUST be a white space after the FIRST "
		pcfa := i.addPrimitiveImmediate(".\"")
		i.code.addInter(pcfa, func(i *Interpreter) {

			if i.readingString { // interpreting from repl
				// get the string from the input stream
				//fmt.Println("DEBUG : reading string from REPL")
				token := i.getNextString()
				fmt.Fprintf(i.writer, "%s", token)
			} else { // interpreting from a compound word
				// read the string from memory
				//fmt.Printf("DEBUG : reading string from Memory, ip = %d, rs = %+v\n", i.ip, i.rs)
				rip := i.ip + 1
				len := i.mem[rip]                    // get string lenth
				var k int                            // rune pointer
				for k = rip + 1; k <= len+rip; k++ { // retrieve string
					fmt.Fprintf(i.writer, "%s", string(rune(i.mem[k])))
				}
				i.ip += len + 1
			}
			i.moveIP()
		})
		// in compile mode (immediate), will write the cfa, length and string
		// in the dictionnary
		i.code.addCompil(pcfa, func(i *Interpreter) {

			// fmt.Printf("DEBUG : Cmode: %v, word: %+v\n", i.compileMode, pcfa)
			token := i.getNextString()
			// fmt.Printf("DEBUG : comile mode, read string : %s\n", token)
			rtok := []rune(token) // group by rune
			if i.Err != nil {
				return
			}
			// fmt.Printf("DEBUG : Cmode: %v, word: %+v\n", i.compileMode, pcfa)
			i.mem = append(i.mem, pcfa, len(rtok))
			// store the token, rune by rune
			for _, r := range rtok {
				i.mem = append(i.mem, int(r))
			}
			// in compile mode + immediate, so just make ip = 0
			// to force reading next token
			i.ip = 0
			// i.moveIP()
		})
	}

	// ( -- addr) addr of the BASE user variable.
	i.code.addInter(i.addPrimitive("base"), func(i *Interpreter) {
		i.Err = i.ds.push(UVBase)
		i.moveIP()
	})

	// enter into intrepretation mode, immediate word
	// rs is unchanged
	{
		pcfa := i.addPrimitiveImmediate("[")
		i.code.addInter(pcfa, func(i *Interpreter) {
			i.Err = errors.New("you cannot call '[' except when already in compile mode")
		})
		i.code.addCompil(pcfa, func(i *Interpreter) {
			//fmt.Printf("DEBUG : before [,  ip:%d and rs:%+v\n", i.ip, i.rs.data)
			//i.ip, i.Err = i.rs.pop() // pop out of the wrapper
			i.compileMode = false
			i.ip = 0 // pop in a normal interpreting state
			//fmt.Printf("DEBUG : after [, now ip:%d and rs:%+v\n", i.ip, i.rs.data)
		})
	}

	// enter into compil mode
	// rs is unchanged
	i.code.addInter(i.addPrimitive("]"), func(i *Interpreter) {
		//fmt.Printf("DEBUG : before ],  ip:%d and rs:%+v\n", i.ip, i.rs.data)
		i.compileMode = true
		i.ip = 0 // trigger next word read
		//fmt.Printf("DEBUG : after ],  ip:%d and rs:%+v\n", i.ip, i.rs.data)
	})

	// make last word immediate
	i.code.addInter(i.addPrimitive("immediate"), func(i *Interpreter) {
		i.words[i.lastNfa].immediate = true
		i.moveIP()
	})

	// ( n -- ) comma, add n in the dictionnary
	i.code.addInter(i.addPrimitive(","), func(i *Interpreter) {
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
	i.code.addInter(i.addPrimitive("allot"), func(i *Interpreter) {
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
	i.code.addInter(i.addPrimitive("here"), func(i *Interpreter) {
		i.Err = i.ds.push(len(i.mem))
		i.moveIP()
	})

	// (n addr -- ) store n at the given address, assume memory is allocated
	i.code.addInter(i.addPrimitive("!"), func(i *Interpreter) {
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
	i.code.addInter(i.addPrimitive("@"), func(i *Interpreter) {
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
	i.code.addInter(i.addPrimitive("cr"), func(i *Interpreter) {
		fmt.Fprintln(i.writer)
		i.moveIP()
	})

	// ( -- ) emit space
	i.code.addInter(i.addPrimitive("space"), func(i *Interpreter) {
		fmt.Fprint(i.writer, " ")
		i.moveIP()
	})

	// ( -- ) emit tab
	i.code.addInter(i.addPrimitive("tab"), func(i *Interpreter) {
		fmt.Fprint(i.writer, "\t")
		i.moveIP()
	})

	// ( rune -- ) emit the provided rune
	i.code.addInter(i.addPrimitive("emit"), func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		fmt.Fprintf(i.writer, "%s", string(rune(n)))
		i.moveIP()
	})

	// ( n -- ) drop top of stack
	i.code.addInter(i.addPrimitive("drop"), func(i *Interpreter) {
		_, i.Err = i.ds.pop()
		i.moveIP()
	})

	// ( n -- n n ) dup to of stack
	i.code.addInter(i.addPrimitive("dup"), func(i *Interpreter) {
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n, n)
		i.moveIP()
	})

	// ( n1 n2 n3 -- n2 n3 n1) rotate the stack
	i.code.addInter(i.addPrimitive("rot"), func(i *Interpreter) {
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
	i.code.addInter(i.addPrimitive("swap"), func(i *Interpreter) {
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
	// implementation should account
	// for the fact that r@, r> and >r are wrappers.
	i.code.addInter(i.addPrimitive("r>"), func(i *Interpreter) {
		var r, top int
		top, _ = i.rs.pop()
		r, i.Err = i.rs.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(r)
		i.rs.push(top)
		i.moveIP()
	})

	// ( r -- ) push r into rs
	// implementation should account
	// for the fact that r@, r> and >r are wrappers.
	i.code.addInter(i.addPrimitive(">r"), func(i *Interpreter) {
		var r, top int
		top, _ = i.rs.pop()
		r, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.rs.push(r)
		i.Err = i.rs.push(top)
		i.moveIP()
	})

	// (-- r) push top of rs to ds
	// implementation should account
	// for the fact that r@, r> and >r are wrappers.
	i.code.addInter(i.addPrimitive("r@"), func(i *Interpreter) {
		var r, top int
		top, _ = i.rs.pop()
		r, i.Err = i.rs.top()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(r)
		i.rs.push(top)
		i.moveIP()
	})

	// ( -- ) forget <word> : forget the specified word
	// and all the following content, whatever the vocabulary.
	i.code.addInter(i.addPrimitive("forget"), func(i *Interpreter) {
		token := i.getNextString()
		if i.Err != nil {
			return
		}
		nfa2forget := i.lookup(token)
		// token not found, do nothing !
		if i.Err != nil {
			return
		}
		i.lastNfa = i.mem[nfa2forget]

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
	i.code.addInter(i.addPrimitive("over"), func(i *Interpreter) {
		var n1, n2 int
		n2, _ = i.ds.pop()
		n1, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.Err = i.ds.push(n1, n2, n1)
		i.moveIP()
	})

}
