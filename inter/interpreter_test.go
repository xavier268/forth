package inter

import (
	"fmt"
	"testing"
)

func TestLookupToken(t *testing.T) {
	var s string
	var err error
	var nfa int
	var w *word

	i := NewInterpreter()
	i.dumpuservars()
	i.dumpwords()
	i.dumpmem()

	s = "+"
	nfa, w, err = i.lookup(s)
	if err != nil {
		t.Fatal(s, "==>", nfa, w, err)
	}

	s = ";"
	nfa, w, err = i.lookup(s)
	if err != nil {
		t.Fatal(s, "==>", nfa, w, err)
	}

	s = "nonexistent"
	nfa, w, err = i.lookup(s)
	if err == nil {
		fmt.Println(s, "==>", nfa, w, err)
		t.Fail()
	}

}
