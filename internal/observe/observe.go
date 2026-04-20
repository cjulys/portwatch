// Package observe provides a lightweight event bus for broadcasting
// port-change events to multiple subscribers within portwatch.
package observe

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Event carries a labelled snapshot of port diffs delivered to subscribers.
type Event struct {
	Topic string
	Diffs []scanner.Diff
}

// Handler is a function that receives an Event.
type Handler func(Event)

// Bus is a simple pub/sub event bus.
type Bus struct {
	mu       sync.RWMutex
	subs     map[string][]Handler
	fallback func(string)
}

// New returns an initialised Bus. fallback is called with a diagnostic message
// when a subscriber panics; if nil, panics are silently recovered.
func New(fallback func(string)) *Bus {
	if fallback == nil {
		fallback = func(string) {}
	}
	return &Bus{
		subs:     make(map[string][]Handler),
		fallback: fallback,
	}
}

// Subscribe registers h to receive events published on topic.
func (b *Bus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], h)
}

// Publish sends e to every handler subscribed to e.Topic.
// Each handler is called synchronously; panics are recovered and forwarded to
// the fallback writer so that a misbehaving subscriber cannot crash the daemon.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subs[e.Topic]))
	copy(handlers, b.subs[e.Topic])
	b.mu.RUnlock()

	for _, h := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					b.fallback("observe: subscriber panic recovered")
				}
			}()
			h(e)
		}()
	}
}

// SubscriberCount returns the number of handlers registered for topic.
func (b *Bus) SubscriberCount(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs[topic])
}
