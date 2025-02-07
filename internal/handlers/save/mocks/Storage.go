// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	multipart "mime/multipart"

	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// GetStoragePath provides a mock function with given fields:
func (_m *Storage) GetStoragePath() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStoragePath")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetStorageType provides a mock function with given fields:
func (_m *Storage) GetStorageType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStorageType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SaveFile provides a mock function with given fields: file, name
func (_m *Storage) SaveFile(file multipart.File, name string) error {
	ret := _m.Called(file, name)

	if len(ret) == 0 {
		panic("no return value specified for SaveFile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(multipart.File, string) error); ok {
		r0 = rf(file, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
