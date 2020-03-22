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
	// BUGGY ? TO DO - fix ??
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

	f(t, "r>", "")
	f(t, "r@ .", " 0")

	f(t, `: XX r> ; `, "")                           // forced return
	f(t, `: XX r> ." not reached" ; XX `, "")        // stop execution by popping stack
	f(t, ": XX r@ ; XX  here - . ", " -1")           // check r@ is pointing to the correct address
	f(t, ": XX r@ noop  ; XX  here - . ", " -2")     // check r@ is pointing to the correct address
	f(t, ": XX r@ noop noop ; XX  here - . ", " -3") // check r@ is pointing to the correct address

	// nested
	f(t, `: XX r> ." test1" ; : YY XX ." test2" ; YY `, "test2")
	f(t, `: XX ." test1" ; : YY XX ." test2" ; YY `, "test1test2")

	// balanced rs tests
	f(t, ` : test >r r> ; 1000  test 2000 .  .`, " 2000 1000")

}

func TestVariable(t *testing.T) {

	t.Skip()

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

	t.Skip()

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
