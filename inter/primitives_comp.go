package inter

// compile special compile mode primitives,
// or just do the default compilation.
func (i *Interpreter) compilePrim(wcfa int, w *word) {

	if i.Err != nil {
		return
	}

	switch w.name {

	case "LITERAL": // get n from stack and compile it (n -- )
		var n int
		n, i.Err = i.ds.pop()
		if i.Err != nil {
			return
		}
		i.compileNum(n)

	default:
		i.alloc(1)
		i.mem[len(i.mem)-1] = wcfa
		return
	}
}

// compile a litteral number
func (i *Interpreter) compileNum(num int) {

	if i.Err != nil {
		return
	}

	nfalitt := i.lookupPrimitive("LITERAL")
	if i.Err != nil {
		panic("LITERAL not defined as primitive ?")
	}
	// write cfa of "literal" and number
	i.alloc(2)
	h := len(i.mem)
	i.mem[h-2], i.mem[h-1] = nfalitt+1, num
	return
}
