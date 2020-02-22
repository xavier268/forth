package inter

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
	i.mem = make([]int, UVEndOfDefinition)
	i.mem[UVHere] = UVEndOfDefinition
}
