package postgres

import (
	"context"
	"fmt"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *DB) *EventRepository {
	return &EventRepository{db: db.DB}
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
