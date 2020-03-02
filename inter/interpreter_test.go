package inter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestConstructInterprter(t *testing.T) {
	i := NewInterpreter()
	if i == nil {
		t.Fatal("Cannot construct interpreter")
	}
	if i.Err != nil {
		t.Fatal("Error constructing interpreter : " + i.Err.Error())
	}
}

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
func TestPrint(t *testing.T) {

	// use at repl level
	f(t, ` ." hello world " `, "hello world ")
	f(t, ` ." hello world" `, "hello world")
	f(t, ` ."    hello world" `, "   hello world") // only the first space is eaten up
	f(t, ` ." hello world" " `, "hello world", true)

	f(t, "DECIMAL 3564 EMIT ", "෬")

	// use inside a definition !
	f(t, ": t .\" hello world\" ;  1 . t", " 1hello world")
	f(t, ": t .\" hello world\" ;   t 1 . ", "hello world 1")

}

func TestOperations(t *testing.T) {

	f(t, "2 3 + ", "")
	f(t, "2 3 + . ", " 5")
	f(t, "2 3 . ", " 3")
	f(t, "2 3 4 . + .", " 4 5")

	f(t, "2 3 - . ", " -1")
	f(t, "3 2 - . ", " 1")

	f(t, "3 2 * . ", " 6")

	f(t, "2 3 SWAP . .  ", " 2 3")
	f(t, "3 DUP +  .  ", " 6")
	f(t, "3 DROP  ", "")
	f(t, "3 DROP . ", "", true)
	f(t, "3 4 DROP . ", " 3")
	f(t, "1 2 OVER . . . ", " 1 2 1")
	f(t, "2 OVER ", "", true)

	f(t, "2 ROT ", "", true)
	f(t, "1 2 ROT ", "", true)
	f(t, "1 2 3 ROT ", "")
	f(t, "1 2 3 ROT . . . ", " 1 3 2")

	f(t, ".", "", true)        // overflown error expected
	f(t, ". 1 . ", "", true)   // overflow, then normal operation
	f(t, " 1 . .", " 1", true) // normal then overflow

}
func TestConstantAndForget(t *testing.T) {
	f(t, "CONSTANT", "", true)
	f(t, "1 CONSTANT", "", true)
	f(t, "55 CONSTANT CC CC . ", " 55")
	f(t, "55 CONSTANT CC : CCC CC CC + . ; CCC", " 110")

	f(t, "4 CONSTANT Q Q . FORGET Q Q . ", " 4", true)
	f(t, "4 CONSTANT Q : R Q ;  FORGET Q ", "")
	f(t, "4 CONSTANT Q : R Q ;  FORGET Q R ", "", true)

}

func TestBuildDoes(t *testing.T) {

	in := `	: DEF [ ." aa"] ." bb" <BUILDS ." cc" DOES> ." dd" ; `
	out := "aa"
	f(t, in, out)

	in += " DEF xx"
	out += "bbcc"
	f(t, in, out)

	in2 := in + " ."
	out += ""
	f(t, in2, out, true) // stack underflow

	in += " xx"
	out += "dd"
	f(t, in, out)

	in += " DROP" // addr should be on stack
	out += ""
	f(t, in, out)

	// use in compound word ...
	in += " : TT xx ;"
	out += ""
	f(t, in, out)
	f(t, in+" .", out, true) // stack underflow expected

	in += " TT DROP"
	out += "dd"
	f(t, in, out)

}

func TestReturnStack(t *testing.T) {
	f(t, "R>", "", true) // stack underflow
	f(t, "R@", "", true) // stack underflow
	f(t, `: XX R> ; `, "")
	f(t, `: XX R> ; XX `, "", true) // stack underflow
	f(t, ": XX R>  ; : YY XX ; YY  HERE - . ", " -1")

	// unbalanced return stack test, implemenation dependent
	f(t, `: XX R>  ; : YY XX ." never displayed " ; YY  HERE - . `, " -19")

	// balanced rs tests
	f(t, ` : test  >R DUP R> ; 1000 2000 test . . .`, " 2000 1000 1000")

}

func TestVariable(t *testing.T) {
	f(t, "VARIABLE v v @ .", " 0")
	f(t, "VARIABLE v 555 v ! v @ . ", " 555")
}

func TestComment(t *testing.T) {

	f(t, "2 3 ( 55 kjhkjh ) + ", "")
	f(t, "2 3 + . ", " 5")
	f(t, "2 ( ; kjhkjh ) 3 . ", " 3")
	f(t, "2 3 ( 33 ) 4 . + .", " 4 5")

	f(t, ": plus + ( ; <- immediate word have no effect ) . ; "+
		": plusplus plus plus ; "+
		"1 2 3 4 plusplus",
		" 7 3")

}
func TestHereAllot(t *testing.T) {

	f(t, "HERE @", "", true)
	f(t, "HERE HERE - . ", " 0")
	f(t, "HERE 1 - @", "")
	f(t, "HERE 3 ALLOT HERE - . ", " -3")
	f(t, "55 , HERE 1 - @ . ", " 55")
	f(t, "2 ALLOT 55  HERE 2 - !  HERE 2 - @ .  ", " 55")

	f(t, ",", "", true) // ds underflow
	f(t, "HERE 1000 , HERE - .", " -1")
	f(t, "666 , 888 , HERE 1 - @ . ", " 888")
	f(t, "666 , 888 , HERE 1 - @ . HERE 2 - @ . ", " 888 666")
}

func TestNoop(t *testing.T) {

	f(t, "NOOP", "")

	f(t, "2 3 + NOOP . ", " 5")
	f(t, "2 NOOP 3 . ", " 3")
	f(t, "NOOP 2 3 4 . + .", " 4 5")

	f(t, ": toto NOOP ; toto", "")
	f(t, ": plus NOOP + ; 3 7 plus .", " 10")
	f(t, ": plus NOOP + . ; 3 7 plus", " 10")

	f(t, ": p1 1 NOOP + NOOP ; : p2 NOOP 2 + ; : p3 p1 p2 ; 5 p3 .", " 8")
	f(t, ": plus + NOOP . ; : plusplus plus NOOP plus ; 1 2 3 4 plusplus",
		" 7 3")

}

func TestDefinition(t *testing.T) {

	f(t, ": toto ; toto", "")
	f(t, ": plus + ; 3 7 plus .", " 10")
	f(t, ": plus + . ; 3 7 plus", " 10")

	f(t, ": p1 1 + ; : p2 2 + ; : p3 p1 p2 ; 5 p3 .", " 8")
	f(t, ": plus + . ; : plusplus plus plus ; 1 2 3 4 plusplus", " 7 3")

	f(t, ": toto [ 3 2 + ] LITERAL ; 1 toto . . ", " 5 1")
	f(t, "5 LITERAL ", "", true)

	f(t, " [ 1 2 + .", "", true)
	f(t, "1 ] 2 + .", "")        // 2 and + are compiled ...
	f(t, " : wrong ] 2 . ;", "") // absurd, but not an error

	f(t, " : t 33 . ; IMMEDIATE : tt t ; ", " 33")
	f(t, " : t 33 . ; ( IMMEDIATE ) : tt t ; ", "")

}

// generic test.
// provide a value for expecterror if you expect an error.
func f(t *testing.T, source, expect string, expecterror ...bool) {
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

func TestTick(t *testing.T) {

	// constant value is stored at pfa +1, check vs tick
	f(t, "55 CONSTANT x ' x 1 + @ .", " 55")

	// variable points to pfa +1, check vs tick
	f(t, "VARIABLE v ' v 1 + v -  . ", " 0")

	// compile mode test
	// pfa of v is compiled into test definition, then checked
	f(t, "VARIABLE v : test ' v 1 + ; test v - . ", " 0")

}
