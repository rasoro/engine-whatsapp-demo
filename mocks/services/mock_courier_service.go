// Code generated by MockGen. DO NOT EDIT.
// Source: services/courier_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCourierService is a mock of CourierService interface.
type MockCourierService struct {
	ctrl     *gomock.Controller
	recorder *MockCourierServiceMockRecorder
}

// MockCourierServiceMockRecorder is the mock recorder for MockCourierService.
type MockCourierServiceMockRecorder struct {
	mock *MockCourierService
}

// NewMockCourierService creates a new mock instance.
func NewMockCourierService(ctrl *gomock.Controller) *MockCourierService {
	mock := &MockCourierService{ctrl: ctrl}
	mock.recorder = &MockCourierServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCourierService) EXPECT() *MockCourierServiceMockRecorder {
	return m.recorder
}

// RedirectMessage mocks base method.
func (m *MockCourierService) RedirectMessage(arg0, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedirectMessage", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RedirectMessage indicates an expected call of RedirectMessage.
func (mr *MockCourierServiceMockRecorder) RedirectMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedirectMessage", reflect.TypeOf((*MockCourierService)(nil).RedirectMessage), arg0, arg1)
}
