package inter

// allocate the number of 0 values on the dictionnary.
// Shift the UVHere pointer to the end of UVHere.
// Return the new here value.
func (i *Interpreter) alloc(n int) {
	i.mem = append(i.mem, make([]int, n)...)
	i.here += n
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
