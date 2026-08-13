package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         uuid.UUID `json:"id"`
	ScriptID   uuid.UUID `json:"script_id"`
	User       string    `json:"user"`
	ScriptPath string    `json:"script_path"`
	Action     string    `json:"action"`
	Time       time.Time `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
}

type CallbackRequest struct {
	User       string `json:"user" validate:"required"`
	ScriptPath string `json:"script" validate:"required"`
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
