// Package inter contains the forth interpreter.
package inter

// Detail memory architecture :
//
//
// Structure of a dictionary entry :
//
// NFA 	->  	->	address used as a key to the word
// LFA == NFA	-> 	int	->	points to previous NFA
// CFA	->	int ->	address of definition word (not for primitives)
// PFA	->	int	-> 	parameters or embedded CFAs (not for primitives)
//
//
