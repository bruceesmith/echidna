// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package observable is a simple generic implementation of the Observer design pattern.

Observables are identified by a string name which must be unique across an application.
Calling the [observable/Observe] function both registers an observable by name and type, but also
registers an observer by providing a notification function that is invoked when the value of the
observable changes. Multiple observers can call [observable/Observe] to register for notifications
concerning an existing observable; in this case, notification functions are called in the order
in which they were registered.

To change the value of an observable, and to notify all observers by invoking their
callback notification function, call the [observable/Set] function.
*/
package observable

import (
	"fmt"
	"reflect"
)

type observable []any

var (
	observables = make(map[string]observable)
)

// Observe either registers a new observable, or adds another observer
// (notification function) to an existing observable
func Observe[T any](name string, cb func(T)) error {
	vtype := reflect.TypeOf(cb).In(0)
	callbacks, existing := observables[name]
	if !existing {
		observables[name] = observable{cb}
	} else {
		evtype := reflect.TypeOf(callbacks[0]).In(0)
		if evtype != vtype {
			return fmt.Errorf("cannot add new callback func(%s) to observable %s (type %s)", vtype.String(), name, evtype.String())
		}
		observables[name] = append(callbacks, cb)
	}
	return nil
}

// Set notifies all observers that the value of an observable has
// changed by calling all registered notification functions
func Set[T any](name string, value T) error {
	vtype := reflect.TypeOf(value)
	callbacks, existing := observables[name]
	if !existing {
		return fmt.Errorf("observable %s has not been registered", name)
	}
	evtype := reflect.TypeOf(callbacks[0]).In(0)
	if vtype != evtype {
		return fmt.Errorf("cannot set %s value into observable %s (type %s)", vtype.String(), name, evtype.String())
	}
	for _, cb := range callbacks {
		callback, ok := cb.(func(T))
		if ok {
			callback(value)
		}
	}
	return nil
}
