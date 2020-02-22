package inter

import (
	"fmt"
	"testing"
)

func TestLookupToken(t *testing.T) {
	var s string
	var err error
	var nfa, opcode int
	var w *word

	i := NewInterpreter()
	//i.dumpuservars()
	//i.dumpwords()
	//i.dumpmem()

	s = "+"
	nfa, opcode, w, err = i.lookup(s)
	if err != nil {
		t.Fatal(s, "==>", nfa, opcode, w, err)
	}

	s = ";"
	nfa, opcode, w, err = i.lookup(s)
	if err != nil {
		t.Fatal(s, "==>", nfa, opcode, w, err)
	}

	s = "nonexistent"
	nfa, opcode, w, err = i.lookup(s)
	if err == nil {
		fmt.Println(s, "==>", nfa, opcode, w, err)
		t.Fail()
	}

}
