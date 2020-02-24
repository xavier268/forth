package inter

// allocate the number of 0 values on the dictionnary.
func (i *Interpreter) alloc(n int) {
	i.mem = append(i.mem, make([]int, n, n)...)
}

// createHeader creates a new header in dictionnary.
// updating words, lastNfa.
// return nfa of created header.
func (i *Interpreter) createHeader(token string) (nfa int) {
	nfa = len(i.mem)
	i.alloc(1)
	i.mem[nfa] = i.lastNfa
	i.words[nfa] = &word{token, false, false}
	i.lastNfa = nfa
	return nfa
}
