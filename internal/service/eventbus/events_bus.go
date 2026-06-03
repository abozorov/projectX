package events

type Bus struct {
	events chan Event
}

func NewBus(buffer int) *Bus {
	return &Bus{
		events: make(chan Event, buffer),
	}
}

func (b *Bus) Publish(event Event) {
	b.events <- event
}

func (b *Bus) Subscribe() <-chan Event {
	return b.events
}
