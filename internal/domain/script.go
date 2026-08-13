package domain

import (
	"time"

	"github.com/google/uuid"
)

type Script struct {
	ID        uuid.UUID `json:"id"`
	Host      string    `json:"host"`
	User      string    `json:"user"`
	Template  string    `json:"template"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateScriptRequest struct {
	Host     string `json:"host" validate:"required,hostname|ip"`
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	Template string `json:"template" validate:"required"`
}

type CreateScriptResponse struct {
	ScriptID   string    `json:"script_id"`
	ScriptPath string    `json:"script_path"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewScript(host, user, template, path string) *Script {
	now := time.Now()
	return &Script{
		ID:        uuid.New(),
		Host:      host,
		User:      user,
		Template:  template,
		Path:      path,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Script) Update() {
	s.UpdatedAt = time.Now()
}
