package domain_test

import (
	"testing"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {
	templates := domain.GetTemplates()

	assert.NotEmpty(t, templates)
	assert.Contains(t, templates, "template1")
	assert.Contains(t, templates, "template2")

	tmpl1 := templates["template1"]
	assert.Equal(t, "template1", tmpl1.Name)
	assert.NotEmpty(t, tmpl1.Content)
	assert.Contains(t, tmpl1.Content, "#!/bin/bash")

	tmpl2 := templates["template2"]
	assert.Equal(t, "template2", tmpl2.Name)
	assert.NotEmpty(t, tmpl2.Content)
	assert.Contains(t, tmpl2.Content, "Disk Space Monitor")
}
