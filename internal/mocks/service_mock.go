// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/service.go -destination=internal/mocks/service_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/lutefd/ai-router-go/internal/models"
	service "github.com/lutefd/ai-router-go/internal/service"
	gomock "go.uber.org/mock/gomock"
)

// MockAIServiceInterface is a mock of AIServiceInterface interface.
type MockAIServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAIServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockAIServiceInterfaceMockRecorder is the mock recorder for MockAIServiceInterface.
type MockAIServiceInterfaceMockRecorder struct {
	mock *MockAIServiceInterface
}

// NewMockAIServiceInterface creates a new mock instance.
func NewMockAIServiceInterface(ctrl *gomock.Controller) *MockAIServiceInterface {
	mock := &MockAIServiceInterface{ctrl: ctrl}
	mock.recorder = &MockAIServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAIServiceInterface) EXPECT() *MockAIServiceInterfaceMockRecorder {
	return m.recorder
}

// GenerateDeepSeekResponse mocks base method.
func (m *MockAIServiceInterface) GenerateDeepSeekResponse(ctx context.Context, model, prompt string, callback func(string)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateDeepSeekResponse", ctx, model, prompt, callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateDeepSeekResponse indicates an expected call of GenerateDeepSeekResponse.
func (mr *MockAIServiceInterfaceMockRecorder) GenerateDeepSeekResponse(ctx, model, prompt, callback any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateDeepSeekResponse", reflect.TypeOf((*MockAIServiceInterface)(nil).GenerateDeepSeekResponse), ctx, model, prompt, callback)
}

// GenerateOpenAIResponse mocks base method.
func (m *MockAIServiceInterface) GenerateOpenAIResponse(ctx context.Context, model, prompt string, callback func(string)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateOpenAIResponse", ctx, model, prompt, callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateOpenAIResponse indicates an expected call of GenerateOpenAIResponse.
func (mr *MockAIServiceInterfaceMockRecorder) GenerateOpenAIResponse(ctx, model, prompt, callback any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateOpenAIResponse", reflect.TypeOf((*MockAIServiceInterface)(nil).GenerateOpenAIResponse), ctx, model, prompt, callback)
}

// GenerateResponse mocks base method.
func (m *MockAIServiceInterface) GenerateResponse(ctx context.Context, model, prompt string, callback func(string)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateResponse", ctx, model, prompt, callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateResponse indicates an expected call of GenerateResponse.
func (mr *MockAIServiceInterfaceMockRecorder) GenerateResponse(ctx, model, prompt, callback any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateResponse", reflect.TypeOf((*MockAIServiceInterface)(nil).GenerateResponse), ctx, model, prompt, callback)
}

// MockAuthServiceInterface is a mock of AuthServiceInterface interface.
type MockAuthServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockAuthServiceInterfaceMockRecorder is the mock recorder for MockAuthServiceInterface.
type MockAuthServiceInterfaceMockRecorder struct {
	mock *MockAuthServiceInterface
}

// NewMockAuthServiceInterface creates a new mock instance.
func NewMockAuthServiceInterface(ctrl *gomock.Controller) *MockAuthServiceInterface {
	mock := &MockAuthServiceInterface{ctrl: ctrl}
	mock.recorder = &MockAuthServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceInterface) EXPECT() *MockAuthServiceInterfaceMockRecorder {
	return m.recorder
}

// AuthenticateUser mocks base method.
func (m *MockAuthServiceInterface) AuthenticateUser(ctx context.Context, email, name, googleID string) (*models.User, *service.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUser", ctx, email, name, googleID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*service.TokenPair)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockAuthServiceInterfaceMockRecorder) AuthenticateUser(ctx, email, name, googleID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockAuthServiceInterface)(nil).AuthenticateUser), ctx, email, name, googleID)
}

// GenerateToken mocks base method.
func (m *MockAuthServiceInterface) GenerateToken(user *models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthServiceInterfaceMockRecorder) GenerateToken(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthServiceInterface)(nil).GenerateToken), user)
}

// GenerateTokenPair mocks base method.
func (m *MockAuthServiceInterface) GenerateTokenPair(user *models.User) (*service.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenPair", user)
	ret0, _ := ret[0].(*service.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateTokenPair indicates an expected call of GenerateTokenPair.
func (mr *MockAuthServiceInterfaceMockRecorder) GenerateTokenPair(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenPair", reflect.TypeOf((*MockAuthServiceInterface)(nil).GenerateTokenPair), user)
}

// RefreshAccessToken mocks base method.
func (m *MockAuthServiceInterface) RefreshAccessToken(refreshToken string) (*service.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshAccessToken", refreshToken)
	ret0, _ := ret[0].(*service.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshAccessToken indicates an expected call of RefreshAccessToken.
func (mr *MockAuthServiceInterfaceMockRecorder) RefreshAccessToken(refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAccessToken", reflect.TypeOf((*MockAuthServiceInterface)(nil).RefreshAccessToken), refreshToken)
}

// ValidateToken mocks base method.
func (m *MockAuthServiceInterface) ValidateToken(tokenString string) (*service.Claims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", tokenString)
	ret0, _ := ret[0].(*service.Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockAuthServiceInterfaceMockRecorder) ValidateToken(tokenString any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuthServiceInterface)(nil).ValidateToken), tokenString)
}

// MockChatServiceInterface is a mock of ChatServiceInterface interface.
type MockChatServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockChatServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockChatServiceInterfaceMockRecorder is the mock recorder for MockChatServiceInterface.
type MockChatServiceInterfaceMockRecorder struct {
	mock *MockChatServiceInterface
}

// NewMockChatServiceInterface creates a new mock instance.
func NewMockChatServiceInterface(ctrl *gomock.Controller) *MockChatServiceInterface {
	mock := &MockChatServiceInterface{ctrl: ctrl}
	mock.recorder = &MockChatServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatServiceInterface) EXPECT() *MockChatServiceInterfaceMockRecorder {
	return m.recorder
}

// CreateChat mocks base method.
func (m *MockChatServiceInterface) CreateChat(ctx context.Context, chat *models.Chat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChat", ctx, chat)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateChat indicates an expected call of CreateChat.
func (mr *MockChatServiceInterfaceMockRecorder) CreateChat(ctx, chat any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChat", reflect.TypeOf((*MockChatServiceInterface)(nil).CreateChat), ctx, chat)
}

// DeleteChat mocks base method.
func (m *MockChatServiceInterface) DeleteChat(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChat", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChat indicates an expected call of DeleteChat.
func (mr *MockChatServiceInterfaceMockRecorder) DeleteChat(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChat", reflect.TypeOf((*MockChatServiceInterface)(nil).DeleteChat), ctx, id)
}

// GetChat mocks base method.
func (m *MockChatServiceInterface) GetChat(ctx context.Context, id string) (*models.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChat", ctx, id)
	ret0, _ := ret[0].(*models.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChat indicates an expected call of GetChat.
func (mr *MockChatServiceInterfaceMockRecorder) GetChat(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChat", reflect.TypeOf((*MockChatServiceInterface)(nil).GetChat), ctx, id)
}

// UpdateChat mocks base method.
func (m *MockChatServiceInterface) UpdateChat(ctx context.Context, chat *models.Chat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChat", ctx, chat)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateChat indicates an expected call of UpdateChat.
func (mr *MockChatServiceInterfaceMockRecorder) UpdateChat(ctx, chat any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChat", reflect.TypeOf((*MockChatServiceInterface)(nil).UpdateChat), ctx, chat)
}

// MockUserServiceInterface is a mock of UserServiceInterface interface.
type MockUserServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockUserServiceInterfaceMockRecorder is the mock recorder for MockUserServiceInterface.
type MockUserServiceInterfaceMockRecorder struct {
	mock *MockUserServiceInterface
}

// NewMockUserServiceInterface creates a new mock instance.
func NewMockUserServiceInterface(ctrl *gomock.Controller) *MockUserServiceInterface {
	mock := &MockUserServiceInterface{ctrl: ctrl}
	mock.recorder = &MockUserServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceInterface) EXPECT() *MockUserServiceInterfaceMockRecorder {
	return m.recorder
}

// GetUsersChatList mocks base method.
func (m *MockUserServiceInterface) GetUsersChatList(ctx context.Context, userID string) ([]*models.UserChat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersChatList", ctx, userID)
	ret0, _ := ret[0].([]*models.UserChat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersChatList indicates an expected call of GetUsersChatList.
func (mr *MockUserServiceInterfaceMockRecorder) GetUsersChatList(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersChatList", reflect.TypeOf((*MockUserServiceInterface)(nil).GetUsersChatList), ctx, userID)
}
