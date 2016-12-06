package cqrs

func NewMockEventHandler() *MockEventHandler {
	return &MockEventHandler{
		make([]Event, 0),
	}
}

type MockEventHandler struct {
	events []Event
}

func (m *MockEventHandler) HandleEvent(event Event) {
	m.events = append(m.events, event)
}
