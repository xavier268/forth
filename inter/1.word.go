package inter

import "fmt"

type word struct {
	name      string
	immediate bool
	smudge    bool
	// TODO
	// store the functons in an array indexed on cfa ?
	// avoid defining functions with interpreter parameter ?
	compil   func()
	inter    func()
	nfa, cfa int
}

// createHeader creates a new header in dictionnary.
// updating words, lastNfa.
// return the created object that was added to the words map.
func (i *Interpreter) createHeader(token string) *word {
	nfa := len(i.mem)
	i.mem = append(i.mem, i.lastNfa)
	w := &word{token, false, false, nil, nil, nfa, nfa + 1}
	i.words[nfa] = w
	i.lastNfa = nfa
	return w
}

// lookup most recent token in dictionnary, using the chain of lfa.
func (i *Interpreter) lookup(token string) (nfa int) {
	return i.lookupFrom(i.lastNfa, token)
}

// lookup only among primitives
func (i *Interpreter) lookupPrimitive(token string) (nfa int) {
	return i.lookupFrom(i.lastPrimitiveNfa, token)
}

// lookup from the lastnfa provided.
func (i *Interpreter) lookupFrom(lastnfa int, token string) (nfa int) {

	// start of search with provied nfa
	nfa = lastnfa
	prevnfa := i.mem[nfa]

	// loop until found or no previous lfa
	for nfa > 0 {
		w := i.words[nfa]
		//fmt.Println("Testing : ", nfa, w)
		if w != nil && w.name == token {
			return nfa
		}

		nfa = prevnfa
		prevnfa = i.mem[nfa]
	}
	i.Err = fmt.Errorf("this token is unknown : %s", token)
	return 0
}

// is the provided word a primitive ?
func (i *Interpreter) isPrimitive(w *word) bool {
	return w.nfa <= i.lastPrimitiveNfa
}
