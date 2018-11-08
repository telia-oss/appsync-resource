// Code generated by MockGen. DO NOT EDIT.
// Source: appsync-client.go

// Package mock_resource is a generated GoMock package.
package mock_resource

import (
	reflect "reflect"

	appsync "github.com/aws/aws-sdk-go/service/appsync"
	gomock "github.com/golang/mock/gomock"
	resource "github.com/telia-oss/appsync-resource"
)

// MockAppsync is a mock of Appsync interface
type MockAppsync struct {
	ctrl     *gomock.Controller
	recorder *MockAppsyncMockRecorder
}

// MockAppsyncMockRecorder is the mock recorder for MockAppsync
type MockAppsyncMockRecorder struct {
	mock *MockAppsync
}

// NewMockAppsync creates a new mock instance
func NewMockAppsync(ctrl *gomock.Controller) *MockAppsync {
	mock := &MockAppsync{ctrl: ctrl}
	mock.recorder = &MockAppsyncMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppsync) EXPECT() *MockAppsyncMockRecorder {
	return m.recorder
}

// GetGraphQLApi mocks base method
func (m *MockAppsync) GetGraphQLApi(apiName string) (*appsync.GraphqlApi, error) {
	ret := m.ctrl.Call(m, "GetGraphQLApi", apiName)
	ret0, _ := ret[0].(*appsync.GraphqlApi)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGraphQLApi indicates an expected call of GetGraphQLApi
func (mr *MockAppsyncMockRecorder) GetGraphQLApi(apiName interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGraphQLApi", reflect.TypeOf((*MockAppsync)(nil).GetGraphQLApi), apiName)
}

// UpdateSchema mocks base method
func (m *MockAppsync) UpdateSchema(apiName string, schema []byte) error {
	ret := m.ctrl.Call(m, "UpdateSchema", apiName, schema)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSchema indicates an expected call of UpdateSchema
func (mr *MockAppsyncMockRecorder) UpdateSchema(apiName, schema interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSchema", reflect.TypeOf((*MockAppsync)(nil).UpdateSchema), apiName, schema)
}

// UpdateResolvers mocks base method
func (m *MockAppsync) UpdateResolvers(apiName string, definitions []resource.ResolverDefinition) error {
	ret := m.ctrl.Call(m, "UpdateResolvers", apiName, definitions)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResolvers indicates an expected call of UpdateResolvers
func (mr *MockAppsyncMockRecorder) UpdateResolvers(apiName, definitions interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateResolvers", reflect.TypeOf((*MockAppsync)(nil).UpdateResolvers), apiName, definitions)
}
