package inter

// PRIMITIVES contains the definitions of the primitives as words.
var PRIMITIVES []word = []word{
	// name, immediate, smudge, primitive
	word{"+", false, false, true},
	word{".", false, false, true},
	word{":", false, false, true},
	word{";", true, false, true},
}

func (i *Interpreter) initPrimitives() {

	//           name, immediate, smudge
	i.addPrimitive("+", false, false)
	i.addPrimitive(".", false, false)
	i.addPrimitive(":", false, false)
	i.addPrimitive(";", true, false)
	i.addPrimitive("LITTERAL", false, false)
}

func (i *Interpreter) addPrimitive(name string, immediate, smudge bool) {
	// save current here value.
	nfa := i.mem[UVHere]

	i.alloc(2)                      // allocate 2 cells
	i.mem[nfa] = len(i.words)       // nfa contains word index as opcode
	i.mem[nfa+1] = i.mem[UVLastNfa] // link to previous name
	i.mem[UVLastNfa] = nfa          // update last header

	// store string and bits separately
	i.words[len(i.words)] =
		// name, immediate, smudge, primitive
		&word{name, immediate, smudge, true}
}
