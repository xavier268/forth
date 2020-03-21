package inter

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetNextString(t *testing.T) {
	ns(t, "a b", "a", "b")
	ns(t, "a         b", "a", "b")
	ns(t, "   a         b", "a", "b")
	ns(t, "   a         b ", "a", "b")
	ns(t, ` ." hello" ." world" `, ".\"", "hello", ".\"", "world")
}

func TestGetNextToken(t *testing.T) {
	nt(t, "1 2 3 ", "1 1", "1 2", "1 3")
	nt(t, "1    2 3 ", "1 1", "1 2", "1 3")
	nt(t, "1   2    3", "1 1", "1 2", "1 3")

	nt(t, "1  ( lkj  )    3", "1 1", "1 3")
	nt(t, "1  ( lkj  )", "1 1")
	nt(t, "1  ( lkj  ", "1 1")
	nt(t, "1  ) lkj  ", "1 1")

}

// ============== test utils ===============
func nt(t *testing.T, s string, out ...string) {
	i := NewInterpreter()
	i.SetReader(strings.NewReader((s)))
	i.Err = nil

	for _, exp := range out {
		r := i.getNextToken()
		got := fmt.Sprintf("%d %d", r.t, r.v)
		if got != exp {
			fmt.Println("Source : ", s)
			fmt.Println("EXPECTED : ", exp)
			fmt.Println("GOT      : ", got)
			t.Fatal()
		}
	}

	// read past the limit
	r := i.getNextToken()
	if r.t != errorT {
		t.Fatal("Error was expected, but we got : ", r)
	}

}

func ns(t *testing.T, s string, tokens ...string) {
	i := NewInterpreter()
	i.SetReader(strings.NewReader((s)))

	for _, exp := range tokens {

		got := i.getNextString()
		if got != exp {
			fmt.Println("Source : ", s)
			fmt.Println("EXPECTED : ", exp)
			fmt.Println("GOT      : ", got)
			t.Fatal()
		}
		if i.Err != nil {
			fmt.Println("Source : ", s)
			fmt.Println("EXPECTED : ", exp)
			fmt.Println("GOT      : ", got)
			t.Fatal("i.Err = ", i.Err)
		}

	}

}
