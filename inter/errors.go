package inter

import "errors"

// ErrStackUnderflow error
var ErrStackUnderflow error = errors.New("stack underflow")

// ErrWordNotFound error
var ErrWordNotFound error = errors.New("unknown word")

// ErrQuit normal exit
var ErrQuit error = errors.New("bye, ... exiting")
