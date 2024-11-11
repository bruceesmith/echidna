// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

// Based on public code from John Arundel, goroutine safety added
//
// URL: https://bitfieldconsulting.com/posts/generic-set

package set

import (
	"fmt"
	"sync"
)

// Set is a generics implementation of the set data type
type Set[E comparable] struct {
	lock   sync.RWMutex
	values map[E]struct{}
}

// New creates a new Set
func New[E comparable](vals ...E) *Set[E] {
	s := Set[E]{
		values: make(map[E]struct{}),
	}
	for _, v := range vals {
		s.values[v] = struct{}{}
	}
	return &s
}

// Add puts a new value into a Set
func (s *Set[E]) Add(vals ...E) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range vals {
		s.values[v] = struct{}{}
	}
}

// Clear removes all values from a Set
func (s *Set[E]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values = make(map[E]struct{})
}

// Contains checks if a value is in the Set
func (s *Set[E]) Contains(v E) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.values[v]
	return ok
}

// Delete remove values(s) from a Set
func (s *Set[E]) Delete(vals ...E) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range vals {
		delete(s.values, v)
	}
}

// Difference returns the set of values that are in s (set A) but not in s2 (set B) ... i.e. A - B
func (s *Set[E]) Difference(s2 *Set[E]) *Set[E] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s2.lock.RLock()
	defer s2.lock.RUnlock()
	result := New[E]()
	for _, v := range s.Members() {
		if !s2.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Disjoint returns true if the intersection of s with another set s2 is empty
func (s *Set[E]) Disjoint(s2 *Set[E]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s2.lock.RLock()
	defer s2.lock.RUnlock()
	return s.Intersection(s2).Empty()
}

// Empty returns true if the Set is empty
func (s *Set[E]) Empty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.values) == 0
}

// Intersection returns the logical intersection of 2 Sets
func (s *Set[E]) Intersection(s2 *Set[E]) *Set[E] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s2.lock.RLock()
	defer s2.lock.RUnlock()
	result := New[E]()
	for _, v := range s.Members() {
		if s2.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Members returns a slice of the values in a Set
func (s *Set[E]) Members() []E {
	s.lock.RLock()
	defer s.lock.RUnlock()
	result := make([]E, 0, len(s.values))
	for v := range s.values {
		result = append(result, v)
	}
	return result
}

// Size returns the number of values in a Set
func (s *Set[E]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.values)
}

// String returns a string representation of the Set members
func (s *Set[E]) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return fmt.Sprintf("%v", s.Members())
}

// Union returns the logical union of 2 Sets
func (s *Set[E]) Union(s2 *Set[E]) *Set[E] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s2.lock.RLock()
	defer s2.lock.RUnlock()
	result := New(s.Members()...)
	result.Add(s2.Members()...)
	return result
}
