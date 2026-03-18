package event

import (
	"log/slog"
	"sync"
)

type Handler func(Event)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *Bus) Publish(eventType string, payload any) {
	b.mu.RLock()
	handlers := b.handlers[eventType]
	b.mu.RUnlock()

	e := Event{Type: eventType, Payload: payload}

	for _, h := range handlers {
		go func(handler Handler) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("event handler panicked",
						"event", eventType,
						"error", r,
					)
				}
			}()
			handler(e)
		}(h)
	}
}
