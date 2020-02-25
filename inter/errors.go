package inter

import (
	"errors"
	"fmt"
)

// ErrStackUnderflow error
var ErrStackUnderflow = errors.New("stack underflow")

// ErrReservedWord error
func ErrReservedWord(token string) error {
	return fmt.Errorf("the token '%s' is reserved for internal use", token)
}

// ErrNotPrimitive error
var ErrNotPrimitive = errors.New("not a valid primitive cfa")

// ErrMissingParent error
var ErrMissingParent = errors.New("missing closing parenthesis")

// ErrWordNotFound error
func ErrWordNotFound(token string) error {
	return errors.New("token '" + token + "' was not found")
}

// ErrUnexpectedEndOfLine error
var ErrUnexpectedEndOfLine = errors.New("a token was expected immediatly after")

// ErrInvalidCfa error
func ErrInvalidCfa(cfa int) error {
	return fmt.Errorf("invalid cfa : %d", cfa)
}

// ErrInvalidAddr error
func ErrInvalidAddr(a int) error {
	return fmt.Errorf("invalid address : %d", a)
}

// ErrQuit normal exit
var ErrQuit = errors.New("bye, ... exiting")

// ErrAbort rend la main à l'utilisateur,
// met les piles à 0, interpreter mode
var ErrAbort = errors.New("abort")

// Abort reset stacks and interpreter
func (i *Interpreter) Abort() {
	i.ds.clear()
	i.rs.clear()
	i.compileMode = false
	// don't override previous error
	if i.Err == nil {
		i.Err = ErrAbort
	}
}
