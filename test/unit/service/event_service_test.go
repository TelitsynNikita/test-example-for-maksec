package service_test

import (
	"context"
	"testing"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestEventService_ProcessCallback_Success_Execute(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockScriptRepo := new(MockScriptRepository)

	svc := service.NewEventService(mockEventRepo, mockScriptRepo)

	scriptID := uuid.New()
	scriptPath := "/opt/script-monitor/scripts/test.sh"

	script := &domain.Script{
		ID:   scriptID,
		Path: scriptPath,
	}

	req := domain.CallbackRequest{
		User:       "root",
		ScriptPath: scriptPath,
		Action:     "execute",
		Time:       "2026-08-13T08:09:36Z",
	}

	mockScriptRepo.On("GetByPath", mock.Anything, scriptPath).Return(script, nil)
	mockEventRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := svc.ProcessCallback(context.Background(), req)

	assert.NoError(t, err)
	mockScriptRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestEventService_ProcessCallback_Success_Open(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockScriptRepo := new(MockScriptRepository)

	svc := service.NewEventService(mockEventRepo, mockScriptRepo)

	scriptID := uuid.New()
	scriptPath := "/opt/script-monitor/scripts/test.sh"

	script := &domain.Script{
		ID:   scriptID,
		Path: scriptPath,
	}

	req := domain.CallbackRequest{
		User:       "root",
		ScriptPath: scriptPath,
		Action:     "open",
		Time:       "2026-08-13T08:09:36Z",
	}

	mockScriptRepo.On("GetByPath", mock.Anything, scriptPath).Return(script, nil)
	mockEventRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := svc.ProcessCallback(context.Background(), req)

	assert.NoError(t, err)
	mockScriptRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestEventService_ProcessCallback_ScriptNotFound(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockScriptRepo := new(MockScriptRepository)

	svc := service.NewEventService(mockEventRepo, mockScriptRepo)

	req := domain.CallbackRequest{
		User:       "root",
		ScriptPath: "/opt/script-monitor/scripts/nonexistent.sh",
		Action:     "execute",
		Time:       "2026-08-13T08:09:36Z",
	}

	mockScriptRepo.On("GetByPath", mock.Anything, req.ScriptPath).Return(nil, domain.ErrScriptNotFound)

	err := svc.ProcessCallback(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrScriptNotFound)
}

func TestEventService_ProcessCallback_InvalidTime(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockScriptRepo := new(MockScriptRepository)

	svc := service.NewEventService(mockEventRepo, mockScriptRepo)

	scriptID := uuid.New()
	scriptPath := "/opt/script-monitor/scripts/test.sh"

	script := &domain.Script{
		ID:   scriptID,
		Path: scriptPath,
	}

	req := domain.CallbackRequest{
		User:       "root",
		ScriptPath: scriptPath,
		Action:     "execute",
		Time:       "invalid-time",
	}

	mockScriptRepo.On("GetByPath", mock.Anything, scriptPath).Return(script, nil)

	err := svc.ProcessCallback(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing time")
}

func TestEventService_ProcessCallback_InvalidAction(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockScriptRepo := new(MockScriptRepository)

	svc := service.NewEventService(mockEventRepo, mockScriptRepo)

	scriptID := uuid.New()
	scriptPath := "/opt/script-monitor/scripts/test.sh"

	script := &domain.Script{
		ID:   scriptID,
		Path: scriptPath,
	}

	req := domain.CallbackRequest{
		User:       "root",
		ScriptPath: scriptPath,
		Action:     "modify",
		Time:       "2026-08-13T08:09:36Z",
	}

	mockScriptRepo.On("GetByPath", mock.Anything, scriptPath).Return(script, nil)

	err := svc.ProcessCallback(context.Background(), req)

	assert.Error(t, err)
}
