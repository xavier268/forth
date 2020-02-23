package inter

// allocate the number of 0 values on the dictionnary.
// Shift the UVHere pointer to the end of UVHere.
// Return the new here value.
func (i *Interpreter) alloc(n int) {
	i.mem = append(i.mem, make([]int, n)...)
	i.here += n
}
