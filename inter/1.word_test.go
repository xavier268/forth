package inter

import (
	"fmt"
	"testing"
)

func TestLookupToken(t *testing.T) {

	//t.Skip()

	var s string
	var err error
	var nfa int

	i := NewInterpreter()

	//i.dump()
	s = "+"
	nfa = i.lookup(s)
	if i.Err != nil {
		t.Fatal(s, "==>", nfa, err)
	}
	// last word
	w, ok := i.words[i.lastNfa]
	if !ok {
		t.Fatalf("Words : %v\nlen(words) = %d\nLast word did not exist", i.words, len(i.words))
	}
	s = w.name
	nfa = i.lookup(s)
	if i.Err != nil {
		t.Fatal(s, "==>", nfa, err)
	}

	s = "___non__existent___"
	nfa = i.lookup(s)
	if i.Err == nil {
		fmt.Println(s, "==>", nfa, err)
		t.FailNow()
	}

}
