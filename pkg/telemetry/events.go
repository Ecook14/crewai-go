package telemetry

import (
	"sync"
	"time"
)

// EventType defines the category of a telemetry event.
type EventType string

const (
	EventAgentStarted   EventType = "agent_started"
	EventAgentThinking  EventType = "agent_thinking"
	EventToolStarted    EventType = "tool_started"
	EventToolFinished   EventType = "tool_finished"
	EventAgentFinished  EventType = "agent_finished"
	EventTaskStarted    EventType = "task_started"
	EventTaskFinished   EventType = "task_finished"
	EventSystemLog      EventType = "system_log"
)

// Event represents a single unit of telemetry data pushed to the dashboard.
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	AgentRole string                 `json:"agent_role,omitempty"`
	Payload   map[string]interface{} `json:"payload"`
}

// EventBus handles the subscription and broadcasting of execution events.
type EventBus struct {
	subscribers []chan Event
	mu          sync.RWMutex
}

var GlobalBus = &EventBus{
	subscribers: make([]chan Event, 0),
}

// Subscribe adds a new listener for events.
func (b *EventBus) Subscribe() chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 100)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// Unsubscribe removes a listener.
func (b *EventBus) Unsubscribe(ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, sub := range b.subscribers {
		if sub == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			close(ch)
			break
		}
	}
}

// Publish broadcasts an event to all active subscribers.
func (b *EventBus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	for _, sub := range b.subscribers {
		// Non-blocking send to avoid hanging the engine if a subscriber is slow
		select {
		case sub <- e:
		default:
		}
	}
}
