package ssh_test

import (
	"context"
	"testing"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/ssh"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	timeout := 5 * time.Second
	client := ssh.NewClient(timeout)

	assert.NotNil(t, client)
}

func TestClient_RunCommand_InvalidHost(t *testing.T) {
	client := ssh.NewClient(2 * time.Second)

	ctx := context.Background()
	_, err := client.RunCommand(ctx, "invalid-host", 22, "root", "password", "echo test")

	assert.Error(t, err)
}
