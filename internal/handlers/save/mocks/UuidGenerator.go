// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UuidGenerator is an autogenerated mock type for the UuidGenerator type
type UuidGenerator struct {
	mock.Mock
}

// GenerateUUID provides a mock function with given fields:
func (_m *UuidGenerator) GenerateUUID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GenerateUUID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewUuidGenerator creates a new instance of UuidGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUuidGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *UuidGenerator {
	mock := &UuidGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
