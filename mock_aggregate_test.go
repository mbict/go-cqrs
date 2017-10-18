// Code generated by mockery v1.0.0
package cqrs

import mock "github.com/stretchr/testify/mock"
import uuid "github.com/satori/go.uuid"

// MockAggregate is an autogenerated mock type for the Aggregate type
type MockAggregate struct {
	mock.Mock
}

// AggregateId provides a mock function with given fields:
func (_m *MockAggregate) AggregateId() uuid.UUID {
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

// AggregateName provides a mock function with given fields:
func (_m *MockAggregate) AggregateName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Apply provides a mock function with given fields: _a0
func (_m *MockAggregate) Apply(_a0 Event) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(Event) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearUncommittedEvents provides a mock function with given fields:
func (_m *MockAggregate) ClearUncommittedEvents() {
	_m.Called()
}

// CurrentVersion provides a mock function with given fields:
func (_m *MockAggregate) CurrentVersion() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetUncommittedEvents provides a mock function with given fields:
func (_m *MockAggregate) GetUncommittedEvents() []Event {
	ret := _m.Called()

	var r0 []Event
	if rf, ok := ret.Get(0).(func() []Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Event)
		}
	}

	return r0
}

// HandleCommand provides a mock function with given fields: _a0
func (_m *MockAggregate) HandleCommand(_a0 Command) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(Command) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IncrementVersion provides a mock function with given fields:
func (_m *MockAggregate) IncrementVersion() {
	_m.Called()
}

// StoreEvent provides a mock function with given fields: _a0
func (_m *MockAggregate) StoreEvent(_a0 Event) {
	_m.Called(_a0)
}

// Version provides a mock function with given fields:
func (_m *MockAggregate) Version() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}