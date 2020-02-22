package inter

import "errors"

// ErrStackUnderflow error
var ErrStackUnderflow error = errors.New("stack underflow")

// ErrWordNotFound error
func ErrWordNotFound(token string) error {
	return errors.New("token '" + token + "' was not found")
}

// ErrQuit normal exit
var ErrQuit error = errors.New("bye, ... exiting")
