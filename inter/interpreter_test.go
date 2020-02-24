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

	//i.dump()

	s = "+"
	nfa = i.lookup(s)
	if i.Err != nil {
		t.Fatal(s, "==>", nfa, err)
	}

	s = ";"
	nfa = i.lookup(s)
	if i.Err != nil {
		t.Fatal(s, "==>", nfa, err)
	}

	s = "nonexistent"
	nfa = i.lookup(s)
	if i.Err == nil {
		fmt.Println(s, "==>", nfa, err)
		t.Fail()
	}

}

func TestOperations(t *testing.T) {

	testInOut(t, "2 3 + ", "")
	testInOut(t, "2 3 + . ", " 5")
	testInOut(t, "2 3 . ", " 3")
	testInOut(t, "2 3 4 . + .", " 4 5")

	testInOut(t, "2 3 - . ", " -1")
	testInOut(t, "3 2 - . ", " 1")

	testInOut(t, ".", "", true)        // overflown error expected
	testInOut(t, ". 1 . ", "", true)   // overflow, then normal operation
	testInOut(t, " 1 . .", " 1", true) // normal then overflow

}
func TestComment(t *testing.T) {

	testInOut(t, "2 3 ( 55 kjhkjh ) + ", "")
	testInOut(t, "2 3 + . ", " 5")
	testInOut(t, "2 ( ; kjhkjh ) 3 . ", " 3")
	testInOut(t, "2 3 ( 33 ) 4 . + .", " 4 5")

	testInOut(t, ": plus + ( ; <- immediate word have no effect ) . ; : plusplus plus plus ; 1 2 3 4 plusplus", " 7 3")

}
func TestVars(t *testing.T) {

	testInOut(t, "HERE @", "", true)
	testInOut(t, "HERE HERE - . ", " 0")
	testInOut(t, "HERE 1 - @", "")
	testInOut(t, "HERE 3 ALLOT HERE - . ", " -3")
	testInOut(t, "55 , HERE 1 - @ . ", " 55")
	testInOut(t, "2 ALLOT 55  HERE 2 - !  HERE 2 - @ .  ", " 55")

}

func TestNoop(t *testing.T) {

	testInOut(t, "NOOP", "")

	testInOut(t, "2 3 + NOOP . ", " 5")
	testInOut(t, "2 NOOP 3 . ", " 3")
	testInOut(t, "NOOP 2 3 4 . + .", " 4 5")

	testInOut(t, ": toto NOOP ; toto", "")
	testInOut(t, ": plus NOOP + ; 3 7 plus .", " 10")
	testInOut(t, ": plus NOOP + . ; 3 7 plus", " 10")

	testInOut(t, ": p1 1 NOOP + NOOP ; : p2 NOOP 2 + ; : p3 p1 p2 ; 5 p3 .", " 8")
	testInOut(t, ": plus + NOOP . ; : plusplus plus NOOP plus ; 1 2 3 4 plusplus", " 7 3")

}

func TestDefinition(t *testing.T) {

	testInOut(t, ": toto ; toto", "")
	testInOut(t, ": plus + ; 3 7 plus .", " 10")
	testInOut(t, ": plus + . ; 3 7 plus", " 10")

	testInOut(t, ": p1 1 + ; : p2 2 + ; : p3 p1 p2 ; 5 p3 .", " 8")
	testInOut(t, ": plus + . ; : plusplus plus plus ; 1 2 3 4 plusplus", " 7 3")

}

// generic test.
// provide a value for expecterror if you expect an error.
func testInOut(t *testing.T, source, expect string, expecterror ...bool) {
	in := strings.NewReader(source)
	out := bytes.NewBuffer(nil)

	i := NewInterpreter().SetReader(in).SetWriter(out)

	i.Run()

	if len(expecterror) != 0 && i.Err == nil {
		fmt.Println("Expected error did not happen : ", i.Err)
		fmt.Printf("SOURCE    <%s>\n", source)
		fmt.Printf("OUTPUT    <%s>\n", string(out.Bytes()))
		fmt.Printf("EXPECTED  <%s>\n", expect)

		t.Fatal("unexpected test result")

	}
	if (len(expecterror) == 0 && i.Err != nil) || string(out.Bytes()) != expect {
		fmt.Println("Unexpected error : ", i.Err)
		fmt.Printf("SOURCE    <%s>\n", source)
		fmt.Printf("OUTPUT    <%s>\n", string(out.Bytes()))
		fmt.Printf("EXPECTED  <%s>\n", expect)

		t.Fatal("unexpected test result")
	}

}
