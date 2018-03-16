// Code generated by mockery v1.0.0
package mocks

import component "github.com/ksonnet/ksonnet/component"
import mock "github.com/stretchr/testify/mock"
import unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// Component is an autogenerated mock type for the Component type
type Component struct {
	mock.Mock
}

// DeleteParam provides a mock function with given fields: path, options
func (_m *Component) DeleteParam(path []string, options component.ParamOptions) error {
	ret := _m.Called(path, options)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, component.ParamOptions) error); ok {
		r0 = rf(path, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields: wantsNamedSpaced
func (_m *Component) Name(wantsNamedSpaced bool) string {
	ret := _m.Called(wantsNamedSpaced)

	var r0 string
	if rf, ok := ret.Get(0).(func(bool) string); ok {
		r0 = rf(wantsNamedSpaced)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Objects provides a mock function with given fields: paramsStr, envName
func (_m *Component) Objects(paramsStr string, envName string) ([]*unstructured.Unstructured, error) {
	ret := _m.Called(paramsStr, envName)

	var r0 []*unstructured.Unstructured
	if rf, ok := ret.Get(0).(func(string, string) []*unstructured.Unstructured); ok {
		r0 = rf(paramsStr, envName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*unstructured.Unstructured)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(paramsStr, envName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Params provides a mock function with given fields: envName
func (_m *Component) Params(envName string) ([]component.NamespaceParameter, error) {
	ret := _m.Called(envName)

	var r0 []component.NamespaceParameter
	if rf, ok := ret.Get(0).(func(string) []component.NamespaceParameter); ok {
		r0 = rf(envName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.NamespaceParameter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(envName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetParam provides a mock function with given fields: path, value, options
func (_m *Component) SetParam(path []string, value interface{}, options component.ParamOptions) error {
	ret := _m.Called(path, value, options)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, interface{}, component.ParamOptions) error); ok {
		r0 = rf(path, value, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Summarize provides a mock function with given fields:
func (_m *Component) Summarize() ([]component.Summary, error) {
	ret := _m.Called()

	var r0 []component.Summary
	if rf, ok := ret.Get(0).(func() []component.Summary); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.Summary)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
