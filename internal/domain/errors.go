package domain

import "errors"

var (
	ErrScriptNotFound      = errors.New("script not found")
	ErrScriptAlreadyExists = errors.New("script already exists")
	ErrTemplateNotFound    = errors.New("template not found")
	ErrInvalidAction       = errors.New("invalid action")
	ErrInvalidHost         = errors.New("invalid host")
	ErrInvalidUser         = errors.New("invalid user")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
