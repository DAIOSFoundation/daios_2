// Code generated by MockGen. DO NOT EDIT.
// Source: bugreport.go

// Package bugreport is a generated GoMock package.
package bugreport

import (
	gomock "github.com/dai/go-ipfs/gxlibs/github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockExample is a mock of Example interface
type MockExample struct {
	ctrl     *gomock.Controller
	recorder *MockExampleMockRecorder
}

// MockExampleMockRecorder is the mock recorder for MockExample
type MockExampleMockRecorder struct {
	mock *MockExample
}

// NewMockExample creates a new mock instance
func NewMockExample(ctrl *gomock.Controller) *MockExample {
	mock := &MockExample{ctrl: ctrl}
	mock.recorder = &MockExampleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExample) EXPECT() *MockExampleMockRecorder {
	return m.recorder
}

// someMethod mocks base method
func (m *MockExample) someMethod(arg0 string) string {
	ret := m.ctrl.Call(m, "someMethod", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// someMethod indicates an expected call of someMethod
func (mr *MockExampleMockRecorder) someMethod(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "someMethod", reflect.TypeOf((*MockExample)(nil).someMethod), arg0)
}