// Code generated by mockery v1.0.0. DO NOT EDIT.

package cqrs

import eventbus "github.com/mbict/go-eventbus"
import mock "github.com/stretchr/testify/mock"
import time "time"
import uuid "github.com/google/uuid"

// MockEventStream is an autogenerated mock type for the EventStream type
type MockEventStream struct {
	mock.Mock
}

// AggregateId provides a mock function with given fields:
func (_m *MockEventStream) AggregateId() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// Error provides a mock function with given fields:
func (_m *MockEventStream) Error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventType provides a mock function with given fields:
func (_m *MockEventStream) EventType() eventbus.EventType {
	ret := _m.Called()

	var r0 eventbus.EventType
	if rf, ok := ret.Get(0).(func() eventbus.EventType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(eventbus.EventType)
	}

	return r0
}

// Next provides a mock function with given fields:
func (_m *MockEventStream) Next() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Scan provides a mock function with given fields: _a0
func (_m *MockEventStream) Scan(_a0 EventData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(EventData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Timestamp provides a mock function with given fields:
func (_m *MockEventStream) Timestamp() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// Version provides a mock function with given fields:
func (_m *MockEventStream) Version() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
