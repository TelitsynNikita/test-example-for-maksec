package repository

import (
	"context"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/google/uuid"
)

type ScriptRepository interface {
	Create(ctx context.Context, script *domain.Script) error
	GetByPath(ctx context.Context, path string) (*domain.Script, error)
	Exists(ctx context.Context, path string) (bool, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Script, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
