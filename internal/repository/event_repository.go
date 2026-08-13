package repository

import (
	"context"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/google/uuid"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByScriptID(ctx context.Context, scriptID uuid.UUID) ([]domain.Event, error)
	GetByTimeRange(ctx context.Context, start, end time.Time) ([]domain.Event, error)
}
