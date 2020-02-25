package inter

import "testing"

func TestBase(t *testing.T) {

	f(t, "BASE @ . ", " 10")
	f(t, "DECIMAL BASE @ . ", " 10")
	f(t, "HEX BASE @ . ", " 10")
	f(t, "HEX 8 8 + . ", " 10")
	f(t, "DECIMAL 8 8 + . ", " 16")

	f(t, "8 BASE ! 10 1 - . ", " 7")

}
