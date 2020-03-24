package inter

import (
	"fmt"
	"io"
	"strconv"
)

// allocate the number of 0 values on the dictionnary.
func (i *Interpreter) alloc(n int) {
	i.mem = append(i.mem, make([]int, n, n)...)
}

// get next token from input stream
/*
func (i *Interpreter) scanNextToken() string {

	if !i.scanner.Scan() {
		// EOF
		i.Err = io.EOF
		return ""
	}
	token := i.scanner.Text()

	return token
}
*/

const (
	compoundT = iota
	numberT
	errorT
)

type scanResult struct {
	v     int    // CFA of token, or number value
	t     byte   // type of value
	token string // raw token in string form
	err   error
}

// get next token in raw string format, no specific processing
func (i *Interpreter) getNextString() string {

	if !i.scanner.Scan() {
		// EOF
		i.Err = io.EOF
		return "EOF"
	}
	return i.scanner.Text()
}

// get next token, in a usable format
func (i *Interpreter) getNextToken() scanResult {

	var r scanResult

	if i.Err != nil {
		r.err = i.Err
		r.t = errorT
		return r
	}

	if !i.scanner.Scan() {
		// EOF
		i.Err = io.EOF
		r.err = i.Err
		r.t = errorT
		return r
	}

	r.token = i.scanner.Text()

	// slurp comments if needed
	if r.token == "(" {
		for r.token[len(r.token)-1:] != ")" {
			if !i.scanner.Scan() {
				r.err = fmt.Errorf("unexpected unclosed comment : %s", r.token)
				r.t = errorT
				return r
			}
			r.token = i.scanner.Text()
		}
		// then tail-recurse,
		// you can have multiple successive comments ...
		return i.getNextToken()
	}

	// try to decode token
	nfa := i.lookup(r.token)
	if i.Err == nil {
		r.v = 1 + nfa
		r.t = compoundT
		r.err = nil
		return r
	}

	// so, token could not be decoded
	// reset error and try numbers ...
	i.Err = nil
	if num, err := strconv.ParseInt(r.token, i.getBase(), 64); err == nil {
		r.v = int(num)
		r.t = numberT
		r.err = nil
		return r
	}

	// Token cannot be understood
	r.t = errorT
	r.err = fmt.Errorf("cannot understand the token : %s", r.token)
	return r

}

// isPrimitive test using the fact that the cfa
// contains a negative pseudo code.
func (i *Interpreter) isPrimitive(w *word) bool {

	if w == nil {
		return false
	}
	return w.nfa <= i.lastPrimitiveNfa
}
