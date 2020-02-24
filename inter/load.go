package inter

import (
	"io"
)

// Load external definitions
func (i *Interpreter) Load(ior io.Reader) {

	old := i.scanner
	i.SetReader(ior)
	i.Run()
	i.scanner = old

}
