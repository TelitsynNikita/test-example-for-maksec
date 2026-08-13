package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ScriptRepository struct {
	db *sqlx.DB
}

func NewScriptRepository(db *sqlx.DB) *ScriptRepository {
	return &ScriptRepository{db: db}
}

func (r *ScriptRepository) Create(ctx context.Context, script *domain.Script) error {
	query := `
		INSERT INTO scripts (id, host, user_name, template, path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		script.ID,
		script.Host,
		script.User,
		script.Template,
		script.Path,
		script.CreatedAt,
		script.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create script: %w", err)
	}

	return nil
}

func (r *ScriptRepository) GetByPath(ctx context.Context, path string) (*domain.Script, error) {
	query := `
		SELECT id, host, user_name, template, path, created_at, updated_at
		FROM scripts
		WHERE path = $1
	`

	var script domain.Script
	err := r.db.GetContext(ctx, &script, query, path)
	if err == sql.ErrNoRows {
		return nil, domain.ErrScriptNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get script by path: %w", err)
	}

	return &script, nil
}

func (r *ScriptRepository) Exists(ctx context.Context, path string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM scripts WHERE path = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, path)
	if err != nil {
		return false, fmt.Errorf("failed to check script existence: %w", err)
	}

	return exists, nil
}

func (r *ScriptRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Script, error) {
	query := `
		SELECT id, host, user_name, template, path, created_at, updated_at
		FROM scripts
		WHERE id = $1
	`

	var script domain.Script
	err := r.db.GetContext(ctx, &script, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrScriptNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get script by id: %w", err)
	}

	return &script, nil
}

func (r *ScriptRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM scripts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete script: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrScriptNotFound
	}

	return nil
}
