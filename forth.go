package main

import (
	"os"

	"github.com/xavier268/forth/inter"
)

func main() {
	i := inter.NewInterpreter()

	// Load files if specified on command line
	for _, f := range os.Args[1:] {
		i.LoadFile(f)
	}

	i.Repl()

}
