package inter

// stack is a basic stack type.
type stack struct {
	data []int
}

// newStack constructor.
func newStack() *stack {
	return new(stack)
}

// Push on stack
func (s *stack) push(x int) {
	s.data = append(s.data, x)
}

// clear stack
func (s *stack) clear() {
	s.data = []int{}
}

// Pop from stack
func (s *stack) pop() (int, error) {
	if len(s.data) == 0 {
		return 0, ErrStackUnderflow
	}
	x := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return x, nil
}

// Top show top of stack
func (s *stack) top() (int, error) {
	if len(s.data) == 0 {
		return 0, ErrStackUnderflow
	}
	return s.data[len(s.data)-1], nil
}
