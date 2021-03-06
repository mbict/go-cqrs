// Code generated by mockery v1.0.0. DO NOT EDIT.

package cqrs

import eventbus "github.com/mbict/go-eventbus"
import mock "github.com/stretchr/testify/mock"

// MockEventFactory is an autogenerated mock type for the EventFactory type
type MockEventFactory struct {
	mock.Mock
}

// MakeEvent provides a mock function with given fields: _a0
func (_m *MockEventFactory) MakeEvent(_a0 eventbus.EventType) EventData {
	ret := _m.Called(_a0)

	var r0 EventData
	if rf, ok := ret.Get(0).(func(eventbus.EventType) EventData); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EventData)
		}
	}

	return r0
}
