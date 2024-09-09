// Copyright © 2024 Bruce Smith <bruceesmith@gmait.com>
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

// Contains checks if a value is in the Set
func (s *Set[E]) Contains(v E) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.values[v]
	return ok
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

// String returns a string representation of the Set members
func (s *Set[E]) String() string {
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
