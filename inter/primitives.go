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

	for k := range PRIMITIVES {
		// write nfa and lfa
		i.mem = append(i.mem,
			k,                // NFA points to index
			i.mem[UVLastNfa], // LFA to previous word
			k)                // for primitives, use the index as the CFA
		// It is garanteed that real cfa will be larger.

		i.mem[UVLastNfa] = len(i.mem) - 3
	}
}
