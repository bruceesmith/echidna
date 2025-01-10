package stack

import "sync"

// element is an entry on the stack
type element[T any] struct {
	value T
	next  *element[T]
}

// Stack is a Go stack implementation using a linked  list
// It is go-routine safe
type Stack[T any] struct {
	top  *element[T]
	mut  *sync.Mutex
	size int
}

// IsEmpty returns true if the stack has no elements
func (s *Stack[T]) IsEmpty() bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.size == 0
}

// Peek returns a copy of the top element off the stack
func (s *Stack[T]) Peek() (value T, ok bool) {
	s.mut.Lock()
	defer s.mut.Unlock()

	var v T
	if s.top == nil {
		return v, false
	}

	v = s.top.value
	return v, true
}

// Pop removes the top element and returns it
func (s *Stack[T]) Pop() (value T, ok bool) {
	s.mut.Lock()
	defer s.mut.Unlock()

	var v T
	if s.top == nil {
		return v, false
	}

	v = s.top.value
	s.top = s.top.next
	s.size--

	return v, true
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(v T) {
	s.mut.Lock()
	defer s.mut.Unlock()

	element := &element[T]{
		value: v,
		next:  s.top,
	}

	s.top = element
	s.size++
}

// Size returns the number of elements on the stack
func (s *Stack[T]) Size() int {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.size
}
