package inter

import (
	"errors"
	"fmt"
)

// ErrStackUnderflow error
var ErrStackUnderflow error = errors.New("stack underflow")

// ErrNotPrimitive error
var ErrNotPrimitive = errors.New("not a valid primitive cfa")

// ErrWordNotFound error
func ErrWordNotFound(token string) error {
	return errors.New("token '" + token + "' was not found")
}

// ErrUnexpectedEndOfLine error
var ErrUnexpectedEndOfLine = errors.New("a token was expected immediatly after")

// ErrInvalidCfa error
func ErrInvalidCfa(cfa int) error {
	return errors.New("invalid cfa : " + fmt.Sprint(cfa))
}

// ErrQuit normal exit
var ErrQuit = errors.New("bye, ... exiting")

// Abort reset stacks and interpreter
func (i *Interpreter) Abort() {
	i.ds.clear()
	i.rs.clear()
	// IP ?
}
