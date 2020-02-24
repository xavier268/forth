package inter

// definition of  User Variables
const (
	UVBase = iota
)

func (i *Interpreter) initUserVars() {

	i.mem = append(i.mem, 10)
	i.setBase(10)

}

// ============ accessors ==================
func (i *Interpreter) setBase(base int) {
	i.mem[UVBase] = base
}

func (i *Interpreter) getBase() int {
	return i.mem[UVBase]
}
