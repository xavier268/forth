package inter

import "fmt"

// User variables addresses with content at the start of memory
const (
	UVHere            = iota // First free address on dictionnary
	UVBase                   // decimal, or hex, or ...
	UVIP                     // Instruction pointer
	UVLastNfa                // TODO, replace with vocabularies, points to the NFA of the last word.
	UVEndOfDefinition        // last pseudo user var, used it to know how many user var there are
)

// Initilize user variables, at the start of the memory
func (i *Interpreter) initUserVar() {
	i.alloc(UVEndOfDefinition)
}

// allocate the number of 0 values on the dictionnary.
// Shift the UVHere pointer to the end of UVHere.
func (i *Interpreter) alloc(n int) {
	i.mem = append(i.mem, make([]int, n)...)
	i.mem[UVHere] = len(i.mem)
}

func (i *Interpreter) dumpuservars() {
	fmt.Println("UVHere   = ", UVHere, ":", i.mem[UVHere])
	fmt.Println("UVBase   = ", UVBase, ":", i.mem[UVBase])
	fmt.Println("UVIP     = ", UVIP, ":", i.mem[UVIP])
	fmt.Println("UVLastNFA= ", UVLastNfa, ":", i.mem[UVLastNfa])
}
