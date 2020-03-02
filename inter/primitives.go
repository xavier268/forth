package inter

func (i *Interpreter) initPrimitives() {

	//           name, immediate
	i.addPrimitive("BYE", false)
	i.addPrimitive("ABORT", false)
	i.addPrimitive("RESET", false)
	i.addPrimitive("INFO", false)
	i.addPrimitive("ALLOT", false)
	i.addPrimitive("!", false)
	i.addPrimitive("HERE", false)
	i.addPrimitive("@", false)
	i.addPrimitive("DUP", false)
	i.addPrimitive("SWAP", false)
	i.addPrimitive("DROP", false)
	i.addPrimitive("ROT", false)  // ( n1 n2 n3 -- n2 n3 n1 )
	i.addPrimitive("OVER", false) // ( a b -- a b a )
	i.addPrimitive(",", false)
	i.addPrimitive("+", false)
	i.addPrimitive("*", false)
	i.addPrimitive("-", false)
	i.addPrimitive("BASE", false) // ( -- addr)
	i.addPrimitive("EMIT", false) // ( char -- ) emit the provided utf8 char
	i.addPrimitive(".", false)
	i.addPrimitive(".\"", false)
	i.addPrimitive("CR", false)
	i.addPrimitive(":", false)
	i.addPrimitive(";", true)
	// i.addPrimitive("CONSTANT", false)
	// i.addPrimitive("$$CONSTANT$$", false) // Internal pseudo keywords
	i.addPrimitive("NOOP", false)
	i.addPrimitive("FORGET", false)
	i.addPrimitive("IMMEDIATE", false)
	i.addPrimitive("[", true)
	i.addPrimitive("]", false)

	i.addPrimitive("R>", false) // ( -- n ) pop RS, and put it on DS
	i.addPrimitive("R@", false) // ( -- n ) just copt top of RS to DS
	i.addPrimitive(">R", false) // ( n -- ) push n on top of RS

	i.addPrimitive("<BUILDS", false)
	i.addPrimitive("DOES>", false)
	i.addPrimitive("$$DOES$$", false) // internal pseudo keyword
	// -----------------------------------------
	// special compile mode behaviour
	i.addPrimitive("LITERAL", false) // compile : (n -- ) comp nuber
	// 									interpr : ( -- n) get number

	// ------------------------------------------
	// TODO string processing

	// TODO vocabularies, vlist,

	// TODO smudge, recursion

	// TODO conditions

	// TODO loops, flow control, ..

	// ------------------------------------------

	// flag last primitive nfa
	i.lastPrimitiveNfa = i.lastNfa

	if i.Err != nil {
		panic(i.Err)
	}
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

// get next token
func (i *Interpreter) scanNextToken() string {

	if !i.scanner.Scan() {
		// EOF
		i.Err = ErrUnexpectedEndOfLine
		return ""
	}
	return i.scanner.Text()
}
