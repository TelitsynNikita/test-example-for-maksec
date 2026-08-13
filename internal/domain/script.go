package domain

import (
	"time"

	"github.com/google/uuid"
)

type Script struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Host      string    `json:"host" db:"host"`
	User      string    `json:"user" db:"user_name"`
	Template  string    `json:"template" db:"template"`
	Path      string    `json:"path" db:"path"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateScriptRequest struct {
	Host     string `json:"host" validate:"required,hostname|ip,min=1,max=255"`
	User     string `json:"user" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=1,max=128"`
	Template string `json:"template" validate:"required,min=1,max=50"`
	Port     int    `json:"port" validate:"min=1,max=65535"`
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
