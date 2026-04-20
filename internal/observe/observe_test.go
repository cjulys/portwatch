package observe_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/observe"
	"github.com/user/portwatch/internal/scanner"
)

func makeDiff(port uint16, dir string) scanner.Diff {
	return scanner.Diff{
		Port:      scanner.Port{Port: port, Protocol: "tcp"},
		Direction: dir,
	}
}

func TestPublishDeliveredToSubscriber(t *testing.T) {
	bus := observe.New(nil)
	var got observe.Event
	bus.Subscribe("ports", func(e observe.Event) { got = e })

	e := observe.Event{Topic: "ports", Diffs: []scanner.Diff{makeDiff(80, "opened")}}
	bus.Publish(e)

	if len(got.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(got.Diffs))
	}
	if got.Diffs[0].Port.Port != 80 {
		t.Errorf("expected port 80, got %d", got.Diffs[0].Port.Port)
	}
}

func TestPublishDeliveredToMultipleSubscribers(t *testing.T) {
	bus := observe.New(nil)
	var mu sync.Mutex
	count := 0
	inc := func(_ observe.Event) {
		mu.Lock()
		count++
		mu.Unlock()
	}
	bus.Subscribe("ports", inc)
	bus.Subscribe("ports", inc)
	bus.Subscribe("ports", inc)

	bus.Publish(observe.Event{Topic: "ports"})

	if count != 3 {
		t.Errorf("expected 3 calls, got %d", count)
	}
}

func TestPublishIgnoresUnrelatedTopics(t *testing.T) {
	bus := observe.New(nil)
	called := false
	bus.Subscribe("other", func(_ observe.Event) { called = true })

	bus.Publish(observe.Event{Topic: "ports"})

	if called {
		t.Error("handler for different topic should not be called")
	}
}

func TestSubscriberPanicIsRecovered(t *testing.T) {
	var msg string
	bus := observe.New(func(m string) { msg = m })
	bus.Subscribe("ports", func(_ observe.Event) { panic("boom") })

	// Must not panic the test.
	bus.Publish(observe.Event{Topic: "ports"})

	if msg == "" {
		t.Error("expected fallback message after subscriber panic")
	}
}

func TestSubscriberCount(t *testing.T) {
	bus := observe.New(nil)
	if bus.SubscriberCount("ports") != 0 {
		t.Fatal("expected 0 subscribers initially")
	}
	bus.Subscribe("ports", func(_ observe.Event) {})
	bus.Subscribe("ports", func(_ observe.Event) {})
	if bus.SubscriberCount("ports") != 2 {
		t.Errorf("expected 2 subscribers, got %d", bus.SubscriberCount("ports"))
	}
}
