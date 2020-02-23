package inter

import "fmt"

// dump
func (i *Interpreter) dump() {
	fmt.Printf("\n%+v\n", i)
}

// dump
func (i *Interpreter) dumpmem() {
	fmt.Println("Memory dump, size =  ", len(i.mem))
	for k, v := range i.mem {
		fmt.Printf("\t%4d: %8d\n", k, v)
	}
}

// dump
func (i *Interpreter) dumpwords() {
	fmt.Println("Words dumps, size = ", len(i.words))
	for k, w := range i.words {
		fmt.Printf("\t%4d:%+v\n", k, w)
	}
}
