package service

import (
	"context"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/repository"
)

type EventService struct {
	eventRepo  repository.EventRepository
	scriptRepo repository.ScriptRepository
}

func NewEventService(eventRepo repository.EventRepository, scriptRepo repository.ScriptRepository) *EventService {
	return &EventService{
		eventRepo:  eventRepo,
		scriptRepo: scriptRepo,
	}
}

func (s *EventService) ProcessCallback(ctx context.Context, req domain.CallbackRequest) error {
	script, err := s.scriptRepo.GetByPath(ctx, req.ScriptPath)
	if err != nil {
		return err
	}

	eventTime, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		return err
	}

	event := domain.NewEvent(script.ID, req.User, req.ScriptPath, req.Action, eventTime)
	return s.eventRepo.Create(ctx, event)
}
