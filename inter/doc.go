// Package inter contains the forth interpreter.
package inter

// Detail memory architecture :
//
//
// Structure of a dictionary entry :
//
// NFA 	->	int	->	a word in the word database, with its flags
// LFA	->	int	-> 	points to the previous NFA
// CFA	->	int ->	address of definition word
// PFA	->	int	-> 	parameters or embedded CFAs
//
//
