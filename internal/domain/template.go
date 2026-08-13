package domain

import (
	_ "embed"
	"strings"
)

//go:embed scripts/templates/template1.sh
var template1Content string

//go:embed scripts/templates/template2.sh
var template2Content string

type ScriptTemplate struct {
	Name        string
	Description string
	Content     string
}

func GetTemplates() map[string]ScriptTemplate {
	return map[string]ScriptTemplate{
		"template1": {
			Name:        "template1",
			Description: "Базовый шаблон мониторинга",
			Content:     strings.TrimSpace(template1Content),
		},
		"template2": {
			Name:        "template2",
			Description: "Шаблон для проверки дискового пространства",
			Content:     strings.TrimSpace(template2Content),
		},
	}
}
