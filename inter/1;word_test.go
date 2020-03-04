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
	s = i.words[len(i.words)-1].name
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
