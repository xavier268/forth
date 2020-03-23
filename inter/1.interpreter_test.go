package inter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestNewInterpreter(t *testing.T) {
	i := NewInterpreter()
	if i == nil {
		t.Fatal("Cannot construct interpreter")
	}
	i.dump()
	//fmt.Println(i)
	if i.Err != nil {
		t.Fatal("Error constructing interpreter : " + i.Err.Error())
	}
}
func TestCompoundWord(t *testing.T) {
	f(t, " : test 333 +  ; 555 test . ", " 888")

	f(t, " 1 2 + . ", " 3")
	f(t, " : test 111 ; test  . ", " 111")
	f(t, " : test 111 222  ; test . . ", " 222 111")

	f(t, " : test 333  . ;  test  ", " 333")
	f(t, " : test 333  + . ;  222 test  ", " 555")

}
func TestPrint1(t *testing.T) {

	// use at repl level
	f(t, ` ." hello world " `, "hello world ")
	f(t, ` ." hello world" `, "hello world")
	f(t, ` ."    hello world" `, "   hello world") // only the first space is eaten up
	f(t, ` ." hello world" " `, "hello world", true)

	f(t, " 3564 emit ", "෬")

	// use inside a definition !
	f(t, ": t .\" hello world\" ;  1 . t", " 1hello world")
	f(t, ": t .\" hello world\" ;   t 1 . ", "hello world 1")
	f(t, ": t 1 .\" hello world\" ;   t  . ", "hello world 1")
	f(t, ": t 1 .\" hello world\" ;   . t  ", "", "data stack underflow")

}
func TestPrint2(t *testing.T) {

	f(t, ": t   .\" aaaa\"  .\" bb\" ;      ", "")
	f(t, ": t   .\" aaaa\"  .\" bb\" ;  t   ", "aaaabb")

}

func TestOperations(t *testing.T) {

	f(t, "2 3 + ", "")
	f(t, "2 3 + . ", " 5")
	f(t, "2 3 . ", " 3")
	f(t, "2 3 4 . + .", " 4 5")

	f(t, "2 3 - . ", " -1")
	f(t, "3 2 - . ", " 1")

	f(t, "3 2 * . ", " 6")

	f(t, "2 3 swap . .  ", " 2 3")
	f(t, "3 dup +  .  ", " 6")
	f(t, "3 drop  ", "")
	f(t, "3 drop . ", "", "UNDERFLOW")
	f(t, "3 4 drop . ", " 3")
	f(t, "1 2 over . . . ", " 1 2 1")
	f(t, "2 over ", "", "UNDERFLOW")

	f(t, "2 rot ", "", "UNDERFLOW")
	f(t, "1 2 rot ", "", "UNDERFLOW")
	f(t, "1 2 3 rot ", "")
	f(t, "1 2 3 rot . . . ", " 1 3 2")

	f(t, ".", "", "UNDERFLOW")
	f(t, ". 1 . ", "", "UNDERFLOW")
	f(t, " 1 . .", " 1", "UNDERFLOW")

}
func TestConstantAndForget(t *testing.T) {

	t.Skip()

	f(t, "CONSTANT", "", true)
	f(t, "1 CONSTANT", "", true)
	f(t, "55 CONSTANT CC CC . ", " 55")
	f(t, "55 CONSTANT CC : CCC CC CC + . ; CCC", " 110")

	f(t, "4 CONSTANT Q Q . FORGET Q Q . ", " 4", true)
	f(t, "4 CONSTANT Q : R Q ;  FORGET Q ", "")
	f(t, "4 CONSTANT Q : R Q ;  FORGET Q R ", "", true)

}

func TestBranch(t *testing.T) {

	f(t, ": t branch [ 0 , ] .\" aaaa\" ; t ", "aaaa") // do nothing
	f(t, " 22 branch ", "", "not in this context")     // fail

	f(t, "33 : t branch [ 2 , ] 55 ; t .", " 33")
	f(t, "33 : t branch [ 0 , ] 55 ; t .", " 55")

	f(t, "0 : t 0branch [ 0 , ] .\" aaaa\" ; t ", "aaaa")           // do nothing
	f(t, "1 : t 0branch [ 0 , ] .\" aaaa\" ; t ", "aaaa")           // do nothing
	f(t, ": t 0branch [ 0 , ] .\" aaaa\" ; t ", "", "ds underflow") // fail
	f(t, "1 0branch ", "", "not in this context")                   // fail

	f(t, "33 : t 0 0branch [ 2 , ] 55 ; t .", " 33")
	f(t, "33 : t 0 0branch [ 0 , ] 55 ; t .", " 55")
	f(t, "33 : t 1 0branch [ 2 , ] 55 ; t .", " 55")
	f(t, "33 : t 1 0branch [ 0 , ] 55 ; t .", " 55")

}
func TestBrackets(t *testing.T) {
	f(t, ": test noop [ 1000  . ] 5500 . ; ", " 1000")
	f(t, ": test noop [ 1000  . ] 5500 . ; test ", " 1000 5500")

	// opening bracket in intrepreted mode is illegal
	f(t, " [ 1 . ] 2 . ", "", "cannot call [ when not in compile mode")

	// closing bracket will move into compile mode ... exit with ;
	f(t, "  ] 2 . ; ", "")
	f(t, "  here ] noop noop ; here swap - . ", " 3") // 2(noop) + 1(;) = 3 cells were written

}

func TestBuildDoes(t *testing.T) {

	t.Skip()

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

	in += " drop" // addr should be on stack
	out += ""
	f(t, in, out)

	// use in compound word ...
	in += " : TT xx ;"
	out += ""
	f(t, in, out)
	f(t, in+" .", out, true) // stack underflow expected

	in += " TT drop"
	out += "dd"
	f(t, in, out)

}

func TestReturnStack(t *testing.T) {

	f(t, "r>", "", "underflow")
	f(t, "r@ .", "", "underflow")

	f(t, `: XX r> ; `, "")              // forced return
	f(t, `: XX r> ." ok" ; XX `, "ok")  // popping stack does not affect current level
	f(t, ": XX r@ ; XX   . ", " 0")     // check r@ is pointing to 0, the repl level
	f(t, ": XX r@ noop  ; XX . ", " 0") // check r@ is pointing to 0, the repl level

	// nested
	f(t, `: XX r> ." test1" ; : YY XX ." test2" ; YY `, "test1")   // popping will cancel second level
	f(t, `: XX ." test1" ; : YY XX ." test2" ; YY `, "test1test2") // no popping - check consistency

	// use rs as temp storage
	f(t, ` : test >r r> ; 1000  test 2000 .  .`, " 2000 1000")

}

func TestVariable(t *testing.T) {

	t.Skip()

	f(t, "VARIABLE v v @ .", " 0")
	f(t, "VARIABLE v 555 v ! v @ . ", " 555")
}

func TestComment(t *testing.T) {

	f(t, "2 3 ( 55 kjhkjh )  + . ", " 5")
	f(t, "2 3 ( 55 kjhkjh ) ( 66 ) + . ", " 5")
	f(t, "2 3 + . ", " 5")
	f(t, "2 ( ; kjhkjh ) 3 . ", " 3")
	f(t, "2 3 ( 33 ) 4 . + .", " 4 5")

	f(t, ": plus + ( ; <- immediate word have no effect ) . ; "+
		": plusplus plus plus ; "+
		"1 2 3 4 plusplus",
		" 7 3")

}
func TestHereAllot(t *testing.T) {

	f(t, "here @", "", true)
	f(t, "here here - . ", " 0")
	f(t, "here 1 - @", "")
	f(t, "here 3 allot here - . ", " -3")
	f(t, "55 , here 1 - @ . ", " 55")
	f(t, "2 allot 55  here 2 - !  here 2 - @ .  ", " 55")

	f(t, ",", "", true) // ds underflow
	f(t, "here 1000 , here - .", " -1")
	f(t, "666 , 888 , here 1 - @ . ", " 888")
	f(t, "666 , 888 , here 1 - @ . here 2 - @ . ", " 888 666")
}

func TestNoop(t *testing.T) {

	f(t, "noop", "")

	f(t, "2 3 + noop . ", " 5")
	f(t, "2 noop 3 . ", " 3")
	f(t, "noop 2 3 4 . + .", " 4 5")

	f(t, ": toto noop ; toto", "")
	f(t, ": plus noop + ; 3 7 plus .", " 10")
	f(t, ": plus noop + . ; 3 7 plus", " 10")

	f(t, ": p1 1 noop + noop ; : p2 noop 2 + ; : p3 p1 p2 ; 5 p3 .", " 8")
	f(t, ": plus + noop . ; : plusplus plus noop plus ; 1 2 3 4 plusplus",
		" 7 3")

}

// ===================================================================

// generic test.
// provide something for expecterror if you expect an error.
// it will be printed if error does not happen.
func f(t *testing.T, source, expect string, expecterror ...interface{}) {
	in := strings.NewReader(source)
	out := bytes.NewBuffer(nil)

	i := NewInterpreter().SetReader(in).SetWriter(out)

	i.Run()

	if len(expecterror) != 0 && i.Err == nil {
		fmt.Println("Expected error did not happen : ", expecterror[0])
		fmt.Printf("SOURCE    <%s>\n", source)
		fmt.Printf("OUTPUT    <%s>\n", string(out.Bytes()))
		fmt.Printf("EXPECTED  <%s>\n", expect)

		t.Fatal("unexpected test result - error is missing")

	}
	if (len(expecterror) == 0 && i.Err != nil) || string(out.Bytes()) != expect {
		fmt.Println("Unexpected error : ", i.Err)
		fmt.Printf("SOURCE    <%s>\n", source)
		fmt.Printf("OUTPUT    <%s>\n", string(out.Bytes()))
		fmt.Printf("EXPECTED  <%s>\n", expect)

		t.Fatal("unexpected test result")
	}

}
