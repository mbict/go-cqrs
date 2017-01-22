package cqrs

type Projection interface {
	//HandlesEvents returns a slice of events this projection can subscribe to
	HandlesEvents() []Event

	//HandleEvent is the event handler entry point
	HandleEvent(event Event) error
}
