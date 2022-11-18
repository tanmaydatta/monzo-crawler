// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Element is an autogenerated mock type for the Element type
type Element struct {
	mock.Mock
}

// GetBaseURL provides a mock function with given fields:
func (_m *Element) GetBaseURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetData provides a mock function with given fields:
func (_m *Element) GetData() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GetType provides a mock function with given fields:
func (_m *Element) GetType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetData provides a mock function with given fields: _a0
func (_m *Element) SetData(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewElement interface {
	mock.TestingT
	Cleanup(func())
}

// NewElement creates a new instance of Element. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewElement(t mockConstructorTestingTNewElement) *Element {
	mock := &Element{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
