// Code generated by mockery v1.0.0
package cqrs

import mock "github.com/stretchr/testify/mock"
import uuid "github.com/satori/go.uuid"

// MockEventStore is an autogenerated mock type for the EventStore type
type MockEventStore struct {
	mock.Mock
}

// LoadStream provides a mock function with given fields: aggregateName, aggregateId
func (_m *MockEventStore) LoadStream(aggregateName string, aggregateId uuid.UUID) (EventStream, error) {
	ret := _m.Called(aggregateName, aggregateId)

	var r0 EventStream
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) EventStream); ok {
		r0 = rf(aggregateName, aggregateId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EventStream)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uuid.UUID) error); ok {
		r1 = rf(aggregateName, aggregateId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteEvent provides a mock function with given fields: _a0, _a1
func (_m *MockEventStore) WriteEvent(_a0 string, _a1 ...Event) error {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...Event) error); ok {
		r0 = rf(_a0, _a1...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}