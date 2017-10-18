// Code generated by mockery v1.0.0
package cqrs

import mock "github.com/stretchr/testify/mock"
import uuid "github.com/satori/go.uuid"

// MockAggregateRepository is an autogenerated mock type for the AggregateRepository type
type MockAggregateRepository struct {
	mock.Mock
}

// Load provides a mock function with given fields: aggregateId
func (_m *MockAggregateRepository) Load(aggregateId uuid.UUID) (Aggregate, error) {
	ret := _m.Called(aggregateId)

	var r0 Aggregate
	if rf, ok := ret.Get(0).(func(uuid.UUID) Aggregate); ok {
		r0 = rf(aggregateId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Aggregate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(aggregateId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: aggregate
func (_m *MockAggregateRepository) Save(aggregate Aggregate) error {
	ret := _m.Called(aggregate)

	var r0 error
	if rf, ok := ret.Get(0).(func(Aggregate) error); ok {
		r0 = rf(aggregate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
