package repository

import (
	"context"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
}
