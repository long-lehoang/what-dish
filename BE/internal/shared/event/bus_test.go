package event

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBus_SubscribeAndPublish(t *testing.T) {
	bus := NewBus()

	var received Event
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe("test.event", func(e Event) {
		received = e
		wg.Done()
	})

	bus.Publish("test.event", "hello")

	wg.Wait()

	assert.Equal(t, "test.event", received.Type)
	assert.Equal(t, "hello", received.Payload)
}

func TestBus_MultipleSubscribers(t *testing.T) {
	bus := NewBus()

	var count atomic.Int32
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		bus.Subscribe("multi.event", func(e Event) {
			count.Add(1)
			wg.Done()
		})
	}

	bus.Publish("multi.event", nil)

	wg.Wait()

	assert.Equal(t, int32(3), count.Load())
}

func TestBus_PublishWithNoSubscribers_NoPanic(t *testing.T) {
	bus := NewBus()
	assert.NotPanics(t, func() {
		bus.Publish("no.subscribers", "data")
	})
}

func TestBus_HandlerPanicDoesNotCrashBus(t *testing.T) {
	bus := NewBus()

	var safeReceived atomic.Bool
	var wg sync.WaitGroup
	wg.Add(2)

	// First handler panics
	bus.Subscribe("panic.event", func(e Event) {
		defer wg.Done()
		panic("boom")
	})

	// Second handler should still run
	bus.Subscribe("panic.event", func(e Event) {
		safeReceived.Store(true)
		wg.Done()
	})

	bus.Publish("panic.event", nil)

	wg.Wait()

	assert.True(t, safeReceived.Load(), "second handler should still execute after first panics")
}

func TestBus_DifferentEventTypes(t *testing.T) {
	bus := NewBus()

	var receivedA, receivedB atomic.Bool
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe("event.a", func(e Event) {
		receivedA.Store(true)
		wg.Done()
	})
	bus.Subscribe("event.b", func(e Event) {
		receivedB.Store(true)
	})

	bus.Publish("event.a", nil)

	wg.Wait()
	// Give a brief moment to confirm event.b handler was NOT called.
	time.Sleep(50 * time.Millisecond)

	assert.True(t, receivedA.Load())
	assert.False(t, receivedB.Load(), "event.b handler should not be called for event.a")
}
