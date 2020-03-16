package inter

import "fmt"

// PrimCode contains the code for primitives,
// both for Interpreted or Compil modes.
type PrimCode struct {
	// The pseudo cfa of the primitive is used as a key.
	// It is always a strictly negative number,
	// to differentiate it from real CFAs
	inter, compil map[int](func(i *Interpreter))
	// default functions to use if the cfa is unknown
	defInter, defCompil func(i *Interpreter)
}

// NewPrimCode constructor, default functions are specified.
// No other functions are added.
func NewPrimCode(defI, defC func(i *Interpreter)) *PrimCode {
	return &PrimCode{
		make(map[int]func(i *Interpreter)),
		make(map[int]func(i *Interpreter)),
		defI, defC}
}

// addInter adds code for interpretation of the provided pseudoCFA.
func (p *PrimCode) addInter(pseudocfa int, inter func(i *Interpreter)) *PrimCode {
	if pseudocfa >= 0 {
		panic("pseudo cfa are expected to be strictly negative, but you provided " + fmt.Sprint(pseudocfa))
	}
	p.inter[pseudocfa] = inter
	return p
}

// addCompil adds code for compil mode of the provided pseudoCFA.
func (p *PrimCode) addCompil(pseudocfa int, compil func(i *Interpreter)) *PrimCode {
	if pseudocfa >= 0 {
		panic("pseudo cfa are expcetd to be strictly negative, you provided " + fmt.Sprint(pseudocfa))
	}
	p.compil[pseudocfa] = compil
	return p
}

// Execute the code, base on the compil mode of the interpreter
func (p *PrimCode) do(i *Interpreter, pcfa int) {
	if i.compileMode {
		f := p.compil[pcfa]
		if f == nil {
			p.defCompil(i)
			return
		}
		f(i)
		return
	}
	f := p.inter[pcfa]
	if f == nil {
		p.defInter(i)
		return
	}
	f(i)
	return
}
