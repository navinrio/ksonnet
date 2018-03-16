// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Code generated by mockery v1.0.0
package mocks

import app "github.com/ksonnet/ksonnet/metadata/app"
import component "github.com/ksonnet/ksonnet/component"
import mock "github.com/stretchr/testify/mock"

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// Component provides a mock function with given fields: ksApp, nsName, componentName
func (_m *Manager) Component(ksApp app.App, nsName string, componentName string) (component.Component, error) {
	ret := _m.Called(ksApp, nsName, componentName)

	var r0 component.Component
	if rf, ok := ret.Get(0).(func(app.App, string, string) component.Component); ok {
		r0 = rf(ksApp, nsName, componentName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(component.Component)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(app.App, string, string) error); ok {
		r1 = rf(ksApp, nsName, componentName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Components provides a mock function with given fields: ns
func (_m *Manager) Components(ns component.Namespace) ([]component.Component, error) {
	ret := _m.Called(ns)

	var r0 []component.Component
	if rf, ok := ret.Get(0).(func(component.Namespace) []component.Component); ok {
		r0 = rf(ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.Component)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(component.Namespace) error); ok {
		r1 = rf(ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NSResolveParams provides a mock function with given fields: ns
func (_m *Manager) NSResolveParams(ns component.Namespace) (string, error) {
	ret := _m.Called(ns)

	var r0 string
	if rf, ok := ret.Get(0).(func(component.Namespace) string); ok {
		r0 = rf(ns)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(component.Namespace) error); ok {
		r1 = rf(ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Namespace provides a mock function with given fields: ksApp, nsName
func (_m *Manager) Namespace(ksApp app.App, nsName string) (component.Namespace, error) {
	ret := _m.Called(ksApp, nsName)

	var r0 component.Namespace
	if rf, ok := ret.Get(0).(func(app.App, string) component.Namespace); ok {
		r0 = rf(ksApp, nsName)
	} else {
		r0 = ret.Get(0).(component.Namespace)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(app.App, string) error); ok {
		r1 = rf(ksApp, nsName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Namespaces provides a mock function with given fields: ksApp, envName
func (_m *Manager) Namespaces(ksApp app.App, envName string) ([]component.Namespace, error) {
	ret := _m.Called(ksApp, envName)

	var r0 []component.Namespace
	if rf, ok := ret.Get(0).(func(app.App, string) []component.Namespace); ok {
		r0 = rf(ksApp, envName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.Namespace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(app.App, string) error); ok {
		r1 = rf(ksApp, envName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
