package inter

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoadBasic(t *testing.T) {

	out := bytes.NewBuffer(nil)
	i := NewInterpreter()
	// create a definition and fill the data stack,
	// in multiple calls with separate compilation
	i.Load(strings.NewReader("55 : plus +    "))
	i.Load(strings.NewReader(" . ;  "))

	// store results
	i.SetWriter(out)
	// Verify dictionnary and datastack have been carried over.
	i.Load(strings.NewReader(" 3 plus "))
	if string(out.Bytes()) != " 58" {
		i.dump()
		t.Fatal("Load did not work correctly")
	}

}
