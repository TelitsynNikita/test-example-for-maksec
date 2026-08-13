package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScriptRepository - мок для ScriptRepository
type MockScriptRepository struct {
	mock.Mock
}

func (m *MockScriptRepository) Create(ctx context.Context, script *domain.Script) error {
	args := m.Called(ctx, script)
	return args.Error(0)
}

func (m *MockScriptRepository) GetByPath(ctx context.Context, path string) (*domain.Script, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Script), args.Error(1)
}

func (m *MockScriptRepository) Exists(ctx context.Context, path string) (bool, error) {
	args := m.Called(ctx, path)
	return args.Bool(0), args.Error(1)
}

// MockSSHClient - мок для SSHClient
type MockSSHClient struct {
	mock.Mock
}

func (m *MockSSHClient) RunCommand(ctx context.Context, host string, port int, user, password, command string) (string, error) {
	args := m.Called(ctx, host, port, user, password, command)
	return args.String(0), args.Error(1)
}

func TestScriptService_CreateScript_Success(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "password",
		Template: "template1",
		Port:     22,
	}

	mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	mockSSH.On("RunCommand", mock.Anything, "127.0.0.1", 22, "root", "password", mock.Anything).Return("", nil).Times(3)

	script, err := svc.CreateScript(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, script)
	assert.Equal(t, req.Host, script.Host)
	assert.Equal(t, req.User, script.User)
	assert.Equal(t, req.Template, script.Template)

	mockRepo.AssertExpectations(t)
	mockSSH.AssertExpectations(t)
}

func TestScriptService_CreateScript_CustomPort(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "password",
		Template: "template1",
		Port:     2222,
	}

	mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	mockSSH.On("RunCommand", mock.Anything, "127.0.0.1", 2222, "root", "password", mock.Anything).Return("", nil).Times(3)

	script, err := svc.CreateScript(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, script)
	assert.Equal(t, req.Host, script.Host)

	mockRepo.AssertExpectations(t)
	mockSSH.AssertExpectations(t)
}

func TestScriptService_CreateScript_ShortPassword(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "ab",
		Template: "template1",
		Port:     22,
	}

	mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	mockSSH.On("RunCommand", mock.Anything, "127.0.0.1", 22, "root", "ab", mock.Anything).Return("", nil).Times(3)

	script, err := svc.CreateScript(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, script)
}

func TestScriptService_CreateScript_EmptyPassword(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "",
		Template: "template1",
		Port:     22,
	}

	script, err := svc.CreateScript(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, script)
	assert.Contains(t, err.Error(), "password is required")
}

func TestScriptService_CreateScript_TemplateNotFound(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "password",
		Template: "template999",
		Port:     22,
	}

	script, err := svc.CreateScript(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, script)
	assert.ErrorIs(t, err, domain.ErrTemplateNotFound)
}

func TestScriptService_CreateScript_AlreadyExists(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "password",
		Template: "template1",
		Port:     22,
	}

	mockRepo.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

	script, err := svc.CreateScript(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, script)
	assert.ErrorIs(t, err, domain.ErrScriptAlreadyExists)
}

func TestScriptService_CreateScript_SSHFailure(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "password",
		Template: "template1",
		Port:     22,
	}

	mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	mockSSH.On("RunCommand", mock.Anything, "127.0.0.1", 22, "root", "password", mock.Anything).
		Return("", errors.New("connection refused")).Times(1)

	script, err := svc.CreateScript(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, script)
	assert.Contains(t, err.Error(), "failed to create remote directory")

	mockRepo.AssertExpectations(t)
	mockSSH.AssertExpectations(t)
}

func TestScriptService_CreateScript_DBFailure_Cleanup(t *testing.T) {
	mockRepo := new(MockScriptRepository)
	mockSSH := new(MockSSHClient)

	svc := service.NewScriptService(mockRepo, mockSSH)

	req := domain.CreateScriptRequest{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "password",
		Template: "template1",
		Port:     2222,
	}

	mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	mockSSH.On("RunCommand", mock.Anything, "127.0.0.1", 2222, "root", "password", mock.Anything).Return("", nil).Times(3)
	mockSSH.On("RunCommand", mock.Anything, "127.0.0.1", 2222, "root", "password", mock.Anything).Return("", nil).Times(1)

	script, err := svc.CreateScript(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, script)
	assert.Contains(t, err.Error(), "failed to save script to database")

	mockRepo.AssertExpectations(t)
	mockSSH.AssertExpectations(t)
}
