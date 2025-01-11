package observable

import (
	"fmt"
	"sync"
)

// Observable represents a value that is handled by an Observer pattern
type Observable struct {
	getter    func() interface{}
	observers map[string]func(func() interface{})
}

var (
	lock  sync.RWMutex
	state = make(map[string]Observable, 0)
)

// DetachObserver removes an Observer from an Observable
func DetachObserver(topic, name string) (err error) {
	lock.Lock()
	defer lock.Unlock()
	t, existing := state[topic]
	if !existing {
		return fmt.Errorf("attempt to detach Observer %s from non-registered Observable %s", name, topic)
	}
	_, existing = t.observers[name]
	if !existing {
		return fmt.Errorf("attempt to detach non-registered Observer %s from Observable %s", name, topic)
	}
	delete(t.observers, name)
	return
}

// DeleteObservable removes all references to an Observable
func DeleteObservable(topic string) (err error) {
	lock.Lock()
	defer lock.Unlock()
	_, existing := state[topic]
	if !existing {
		return fmt.Errorf("attempt to detach non-registered Observable %s", topic)
	}
	delete(state, topic)
	return
}

// Notify calls the Update function for every registered Observer, passing the
// Subject's registered Get function to each Observer. The Oberver will use this Get
// function to fetch the Observable's value
func Notify(topic string) (err error) {
	lock.RLock()
	defer lock.RUnlock()
	t, existing := state[topic]
	if !existing {
		return fmt.Errorf("the pub/sub topic %s is not registered", topic)
	}
	for _, updater := range t.observers {
		updater(t.getter)
	}
	return
}

// RegisterObserver registers a new Observer. If the referenced Observable
// has not previously been registered, this function registers it in anticipation that
// a Subject in a separate goroutine will subsequently register
func RegisterObserver(topic, name string, update func(func() interface{})) (err error) {
	lock.Lock()
	defer lock.Unlock()
	t, existing := state[topic]
	if existing {
		if t.observers == nil {
			return fmt.Errorf("the existing pub/sub topic %s has a nil Observers list", topic)
		}
		t.observers[name] = update
	} else {
		t = Observable{
			getter:    nil,
			observers: map[string]func(func() interface{}){name: update},
		}
		state[topic] = t
	}
	return
}

// RegisterSubject registers a new Subject. If the referenced Observable is
// has previously been registered, this function simply saves its getter
func RegisterSubject(topic string, get func() interface{}) (err error) {
	lock.Lock()
	defer lock.Unlock()
	t, existing := state[topic]
	if existing {
		if t.observers == nil {
			return fmt.Errorf("the existing pub/sub topic %s has a nil observers list", topic)
		}
		t.getter = get
	} else {
		t = Observable{
			getter:    get,
			observers: make(map[string]func(func() interface{}), 0),
		}
		state[topic] = t
	}
	return
}
