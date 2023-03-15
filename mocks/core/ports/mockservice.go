// Code generated by MockGen. DO NOT EDIT.
// Source: gophermart/internal/core/ports (interfaces: UserService,OrderService)

// Package ports is a generated GoMock package.
package ports

import (
	context "context"
	domain "gophermart/internal/core/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// AuthenticateUser mocks base method.
func (m *MockUserService) AuthenticateUser(arg0 context.Context, arg1 string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUser", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockUserServiceMockRecorder) AuthenticateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockUserService)(nil).AuthenticateUser), arg0, arg1)
}

// LoginUser mocks base method.
func (m *MockUserService) LoginUser(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserServiceMockRecorder) LoginUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserService)(nil).LoginUser), arg0, arg1, arg2)
}

// RegisterUser mocks base method.
func (m *MockUserService) RegisterUser(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserServiceMockRecorder) RegisterUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserService)(nil).RegisterUser), arg0, arg1, arg2)
}

// MockOrderService is a mock of OrderService interface.
type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

// MockOrderServiceMockRecorder is the mock recorder for MockOrderService.
type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

// NewMockOrderService creates a new mock instance.
func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockOrderService) AddOrder(arg0 context.Context, arg1 *domain.User, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockOrderServiceMockRecorder) AddOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockOrderService)(nil).AddOrder), arg0, arg1, arg2)
}

// GetAllOrders mocks base method.
func (m *MockOrderService) GetAllOrders(arg0 context.Context, arg1 *domain.User) ([]domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOrders", arg0, arg1)
	ret0, _ := ret[0].([]domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOrders indicates an expected call of GetAllOrders.
func (mr *MockOrderServiceMockRecorder) GetAllOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOrders", reflect.TypeOf((*MockOrderService)(nil).GetAllOrders), arg0, arg1)
}

// GetAllWithdrawals mocks base method.
func (m *MockOrderService) GetAllWithdrawals(arg0 context.Context, arg1 *domain.User) ([]domain.Withdrawn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWithdrawals", arg0, arg1)
	ret0, _ := ret[0].([]domain.Withdrawn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWithdrawals indicates an expected call of GetAllWithdrawals.
func (mr *MockOrderServiceMockRecorder) GetAllWithdrawals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWithdrawals", reflect.TypeOf((*MockOrderService)(nil).GetAllWithdrawals), arg0, arg1)
}

// GetUserBalance mocks base method.
func (m *MockOrderService) GetUserBalance(arg0 context.Context, arg1 *domain.User) (domain.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", arg0, arg1)
	ret0, _ := ret[0].(domain.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockOrderServiceMockRecorder) GetUserBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockOrderService)(nil).GetUserBalance), arg0, arg1)
}

// Withdraw mocks base method.
func (m *MockOrderService) Withdraw(arg0 context.Context, arg1 string, arg2 float64, arg3 *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockOrderServiceMockRecorder) Withdraw(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockOrderService)(nil).Withdraw), arg0, arg1, arg2, arg3)
}
