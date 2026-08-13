package repository

import (
	"context"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/google/uuid"
)

type ScriptRepository interface {
	Create(ctx context.Context, script *domain.Script) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Script, error)
	GetByPath(ctx context.Context, path string) (*domain.Script, error)
	GetByHost(ctx context.Context, host string) ([]domain.Script, error)
	Update(ctx context.Context, script *domain.Script) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByPath(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByScriptID(ctx context.Context, scriptID uuid.UUID) ([]domain.Event, error)
	GetByTimeRange(ctx context.Context, start, end time.Time) ([]domain.Event, error)
	GetLatest(ctx context.Context, limit int) ([]domain.Event, error)
	CountByScriptID(ctx context.Context, scriptID uuid.UUID) (int64, error)
}
