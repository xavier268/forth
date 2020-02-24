package inter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestLookupToken(t *testing.T) {

	//t.Skip()

	var s string
	var err error
	var nfa int

	i := NewInterpreter()

	i.dump()

	s = "+"
	nfa, err = i.lookup(s)
	if err != nil {
		t.Fatal(s, "==>", nfa, err)
	}

	s = ";"
	nfa, err = i.lookup(s)
	if err != nil {
		t.Fatal(s, "==>", nfa, err)
	}

	s = "nonexistent"
	nfa, err = i.lookup(s)
	if err == nil {
		fmt.Println(s, "==>", nfa, err)
		t.Fail()
	}

}

func TestOperations(t *testing.T) {

	testInOut(t, "2 3 + ", "")
	testInOut(t, "2 3 + . ", " 5")
	testInOut(t, "2 3 . ", " 3")
	testInOut(t, ".", "") // overflow

	testInOut(t, ": toto ; toto", "")
	testInOut(t, ": plus + ; 3 7 plus .", " 10")

	// This one is buggy - multiple statement in definition does not work
	testInOut(t, ": plus + . ; 3 7 plus", " 10")

}

func testInOut(t *testing.T, source, expect string) {
	in := strings.NewReader(source)
	out := bytes.NewBuffer(nil)

	i := NewInterpreter().SetReader(in).SetWriter(out)

	i.Run()
	if string(out.Bytes()) != expect {
		fmt.Printf("SOURCE    <%s>\n", source)
		fmt.Printf("OUTPUT    <%s>\n", string(out.Bytes()))
		fmt.Printf("EXPECTED  <%s>\n", expect)

		t.Fatal("unexpected test result")
	}

}
