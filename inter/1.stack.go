package inter

import "fmt"

// stack is a basic stack type.
type stack struct {
	data []int
	// overflow limit
	over int
	// error messages
	errOver, errUnder error
}

// newStack constructor.
func newStack() *stack {
	s := new(stack)
	s.over = 1000
	s.errOver = fmt.Errorf("stack overflow ( limit %d cells )", s.over)
	s.errUnder = fmt.Errorf("stack underflow ")
	return s
}

// Push on stack
func (s *stack) push(x ...int) error {
	s.data = append(s.data, x...)
	if len(s.data) > s.over {
		return s.errOver
	}
	return nil
}

// empty test if stack is empty
func (s *stack) empty() bool {
	return len(s.data) == 0
}

// clear stack
func (s *stack) clear() {
	s.data = []int{}
}

// Pop from stack
func (s *stack) pop() (int, error) {
	if len(s.data) == 0 {
		return 0, s.errUnder
	}
	x := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return x, nil
}

// Top show top of stack
func (s *stack) top() (int, error) {
	if len(s.data) == 0 {
		return 0, s.errUnder
	}
	return s.data[len(s.data)-1], nil
}
