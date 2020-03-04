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
	primitiveT = iota
	compoundT
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

// get next token, is a usable form
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

	// slurp comments
	if r.token == "(" {
		for r.token[len(r.token)-1:] != ")" {
			if !i.scanner.Scan() {
				r.err = fmt.Errorf("unexpected unclosed comment : %s", r.token)
				r.t = errorT
				return r
			}
			r.token += i.scanner.Text()
		}
		// load new non-comment token
		if !i.scanner.Scan() {
			// EOF
			i.Err = io.EOF
			r.err = i.Err
			r.t = errorT
			return r
		}
		r.token = i.scanner.Text()
	}

	// test EOF
	if i.Err == io.EOF {
		r.err = io.EOF
		r.t = errorT
		return r
	}

	// identify the nature of the token
	cfa := 1 + i.lookup(r.token)
	if i.Err == nil { //  found !
		if cfa <= 1+i.lastPrimitiveNfa {
			// primitive
			r.v = cfa
			r.t = primitiveT
			return r
		}
		// not primitive
		r.v = cfa
		r.t = compoundT
		return r

	}
	// not found !
	// reset token not found error
	i.Err = nil
	if num, err := strconv.ParseInt(r.token, i.getBase(), 64); err == nil {
		r.v = int(num)
		r.t = numberT
		return r
	}

	// Token cannot be understood
	r.t = errorT
	r.err = fmt.Errorf("cannot understand the token : %s", r.token)
	return r

}
