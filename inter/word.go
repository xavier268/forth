package inter

type word struct {
	name      string
	immediate bool
	smudge    bool
}

// lookup most recent token in dictionnary, using the chain of lfa.
func (i *Interpreter) lookup(token string) (nfa int, err error) {
	return i.lookupFrom(i.lastNfa, token)
}

// lookup only among primitives
func (i *Interpreter) lookupPrimitive(token string) (nfa int, err error) {
	return i.lookupFrom(i.lastPrimitiveNfa, token)
}

// lookup from the lastnfa provided.
func (i *Interpreter) lookupFrom(lastnfa int, token string) (nfa int, err error) {

	// start of search with provied nfa
	nfa = lastnfa
	prevnfa := i.mem[nfa]

	// loop until found or no previous lfa
	for nfa > 0 {
		w := i.words[nfa]
		//fmt.Println("Testing : ", nfa, w)
		if w != nil && w.name == token {
			return nfa, nil
		}

		nfa = prevnfa
		prevnfa = i.mem[nfa]
	}
	return 0, ErrWordNotFound(token)
}
