// Code generated by MockGen. DO NOT EDIT.
// Source: rates_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "get-rates-usdt-grpc-service/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRatesService is a mock of RatesService interface.
type MockRatesService struct {
	ctrl     *gomock.Controller
	recorder *MockRatesServiceMockRecorder
}

// MockRatesServiceMockRecorder is the mock recorder for MockRatesService.
type MockRatesServiceMockRecorder struct {
	mock *MockRatesService
}

// NewMockRatesService creates a new mock instance.
func NewMockRatesService(ctrl *gomock.Controller) *MockRatesService {
	mock := &MockRatesService{ctrl: ctrl}
	mock.recorder = &MockRatesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRatesService) EXPECT() *MockRatesServiceMockRecorder {
	return m.recorder
}

// GetCurrentRate mocks base method.
func (m *MockRatesService) GetCurrentRate(ctx context.Context) (*models.Rate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentRate", ctx)
	ret0, _ := ret[0].(*models.Rate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentRate indicates an expected call of GetCurrentRate.
func (mr *MockRatesServiceMockRecorder) GetCurrentRate(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentRate", reflect.TypeOf((*MockRatesService)(nil).GetCurrentRate), ctx)
}
