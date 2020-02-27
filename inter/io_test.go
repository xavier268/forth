package inter

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestScanFunc(t *testing.T) {

	ts(t, "this .\" ttt    jhg\"jjg hjg",
		"this",
		".\"",
		"ttt    jhg",
		"jjg",
		"hjg")

	ts(t, "this .\" ttt    jhg\"",
		"this",
		".\"",
		"ttt    jhg")

}

func ts(t *testing.T, s string, res ...string) {

	i := NewInterpreter()

	sc := bufio.NewScanner(bytes.NewBuffer([]byte(s)))
	sc.Split(i.newSplitFunction())

	fmt.Println()

	for i := range res {
		if !sc.Scan() {
			t.Fatal("no more token came too early ?")
		}
		if token := sc.Text(); token != res[i] {
			t.Fatal("Expected token :", res[i], " Got :", token)
		}
		fmt.Println(res[i])
	}
	if sc.Scan() {
		t.Fatal("Remaining token when it should be finished : ", sc.Text())
	}

}
