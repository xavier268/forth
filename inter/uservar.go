package inter

// definition of  User Variables
const (
	UVBase = iota
	UVEndMarker
)

func (i *Interpreter) initUserVars() {

	i.mem = append(i.mem, make([]int, UVEndMarker, UVEndMarker)...)
	i.setBase(10)

}

// ============ accessors ==================
func (i *Interpreter) setBase(base int) {
	i.mem[UVBase] = base
}

func (i *Interpreter) getBase() int {
	return i.mem[UVBase]
}
