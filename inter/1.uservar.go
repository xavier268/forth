package inter

// definition of  User Variables
const (
	UVStartMarker = iota // do not use, to avoid zero addressing
	UVBase
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
