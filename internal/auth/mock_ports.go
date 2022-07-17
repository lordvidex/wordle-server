// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/auth/ports.go

// Package auth is a generated GoMock package.
package auth

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	game "github.com/lordvidex/wordle-wf/internal/game"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(name, email, password string) (*game.Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name, email, password)
	ret0, _ := ret[0].(*game.Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(name, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), name, email, password)
}

// FindByEmail mocks base method.
func (m *MockRepository) FindByEmail(email string) (*game.Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", email)
	ret0, _ := ret[0].(*game.Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockRepositoryMockRecorder) FindByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockRepository)(nil).FindByEmail), email)
}

// FindByID mocks base method.
func (m *MockRepository) FindByID(id uuid.UUID) (*game.Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*game.Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRepository)(nil).FindByID), id)
}

// MockPasswordChecker is a mock of PasswordChecker interface.
type MockPasswordChecker struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordCheckerMockRecorder
}

// MockPasswordCheckerMockRecorder is the mock recorder for MockPasswordChecker.
type MockPasswordCheckerMockRecorder struct {
	mock *MockPasswordChecker
}

// NewMockPasswordChecker creates a new mock instance.
func NewMockPasswordChecker(ctrl *gomock.Controller) *MockPasswordChecker {
	mock := &MockPasswordChecker{ctrl: ctrl}
	mock.recorder = &MockPasswordCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordChecker) EXPECT() *MockPasswordCheckerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockPasswordChecker) Check(password, hash string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", password, hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockPasswordCheckerMockRecorder) Check(password, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockPasswordChecker)(nil).Check), password, hash)
}

// MockTokenHelper is a mock of TokenHelper interface.
type MockTokenHelper struct {
	ctrl     *gomock.Controller
	recorder *MockTokenHelperMockRecorder
}

// MockTokenHelperMockRecorder is the mock recorder for MockTokenHelper.
type MockTokenHelperMockRecorder struct {
	mock *MockTokenHelper
}

// NewMockTokenHelper creates a new mock instance.
func NewMockTokenHelper(ctrl *gomock.Controller) *MockTokenHelper {
	mock := &MockTokenHelper{ctrl: ctrl}
	mock.recorder = &MockTokenHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenHelper) EXPECT() *MockTokenHelperMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockTokenHelper) Decode(token Token) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", token)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode.
func (mr *MockTokenHelperMockRecorder) Decode(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockTokenHelper)(nil).Decode), token)
}

// Generate mocks base method.
func (m *MockTokenHelper) Generate(payload interface{}) (Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", payload)
	ret0, _ := ret[0].(Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockTokenHelperMockRecorder) Generate(payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockTokenHelper)(nil).Generate), payload)
}
