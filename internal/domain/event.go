package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ScriptID   uuid.UUID `json:"script_id" db:"script_id"`
	User       string    `json:"user" db:"user_name"`
	ScriptPath string    `json:"script_path" db:"script_path"`
	Action     string    `json:"action" db:"action"`
	Time       time.Time `json:"time" db:"event_time"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type CallbackRequest struct {
	User       string `json:"user" validate:"required,min=1,max=100"`
	ScriptPath string `json:"script" validate:"required,min=1,max=4096"`
	Action     string `json:"action" validate:"required,oneof=execute modify"`
	Time       string `json:"time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type CallbackResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewEvent(scriptID uuid.UUID, user, scriptPath, action string, eventTime time.Time) *Event {
	return &Event{
		ID:         uuid.New(),
		ScriptID:   scriptID,
		User:       user,
		ScriptPath: scriptPath,
		Action:     action,
		Time:       eventTime,
		CreatedAt:  time.Now(),
	}
}
