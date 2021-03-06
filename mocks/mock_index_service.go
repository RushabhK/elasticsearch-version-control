// Code generated by MockGen. DO NOT EDIT.
// Source: service/index_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIndexService is a mock of IndexService interface
type MockIndexService struct {
	ctrl     *gomock.Controller
	recorder *MockIndexServiceMockRecorder
}

// MockIndexServiceMockRecorder is the mock recorder for MockIndexService
type MockIndexServiceMockRecorder struct {
	mock *MockIndexService
}

// NewMockIndexService creates a new mock instance
func NewMockIndexService(ctrl *gomock.Controller) *MockIndexService {
	mock := &MockIndexService{ctrl: ctrl}
	mock.recorder = &MockIndexServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIndexService) EXPECT() *MockIndexServiceMockRecorder {
	return m.recorder
}

// CreateIndex mocks base method
func (m *MockIndexService) CreateIndex(indexName, configuration string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndex", indexName, configuration)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIndex indicates an expected call of CreateIndex
func (mr *MockIndexServiceMockRecorder) CreateIndex(indexName, configuration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndex", reflect.TypeOf((*MockIndexService)(nil).CreateIndex), indexName, configuration)
}

// ReIndex mocks base method
func (m *MockIndexService) ReIndex(sourceIndex, targetIndex, script string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReIndex", sourceIndex, targetIndex, script)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReIndex indicates an expected call of ReIndex
func (mr *MockIndexServiceMockRecorder) ReIndex(sourceIndex, targetIndex, script interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReIndex", reflect.TypeOf((*MockIndexService)(nil).ReIndex), sourceIndex, targetIndex, script)
}

// DeleteIndex mocks base method
func (m *MockIndexService) DeleteIndex(indexName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", indexName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIndex indicates an expected call of DeleteIndex
func (mr *MockIndexServiceMockRecorder) DeleteIndex(indexName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockIndexService)(nil).DeleteIndex), indexName)
}

// GetDocumentsCount mocks base method
func (m *MockIndexService) GetDocumentsCount(indexName string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDocumentsCount", indexName)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDocumentsCount indicates an expected call of GetDocumentsCount
func (mr *MockIndexServiceMockRecorder) GetDocumentsCount(indexName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDocumentsCount", reflect.TypeOf((*MockIndexService)(nil).GetDocumentsCount), indexName)
}
