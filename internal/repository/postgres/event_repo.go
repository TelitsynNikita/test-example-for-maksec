package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events (id, script_id, user_name, script_path, action, event_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.ScriptID,
		event.User,
		event.ScriptPath,
		event.Action,
		event.Time,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (r *EventRepository) GetByScriptID(ctx context.Context, scriptID uuid.UUID) ([]domain.Event, error) {
	query := `
		SELECT id, script_id, user_name, script_path, action, event_time, created_at
		FROM events
		WHERE script_id = $1
		ORDER BY event_time DESC
	`

	var events []domain.Event
	err := r.db.SelectContext(ctx, &events, query, scriptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by script id: %w", err)
	}

	return events, nil
}

func (r *EventRepository) GetByTimeRange(ctx context.Context, start, end time.Time) ([]domain.Event, error) {
	query := `
		SELECT id, script_id, user_name, script_path, action, event_time, created_at
		FROM events
		WHERE event_time BETWEEN $1 AND $2
		ORDER BY event_time DESC
	`

	var events []domain.Event
	err := r.db.SelectContext(ctx, &events, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by time range: %w", err)
	}

	return events, nil
}
