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
	client := ssh.NewClient(ssh.Config{
		Timeout:               timeout,
		StrictHostKeyChecking: false,
		KnownHostsFile:        "",
	})

	assert.NotNil(t, client)
}

func TestClient_RunCommand_InvalidHost(t *testing.T) {
	client := ssh.NewClient(ssh.Config{
		Timeout:               2 * time.Second,
		StrictHostKeyChecking: false,
		KnownHostsFile:        "",
	})

	ctx := context.Background()
	_, err := client.RunCommand(ctx, "invalid-host", 22, "root", "password", "echo test")

	assert.Error(t, err)
}

func TestClient_RunCommand_WithStrictHostKeyChecking(t *testing.T) {
	client := ssh.NewClient(ssh.Config{
		Timeout:               2 * time.Second,
		StrictHostKeyChecking: true,
		KnownHostsFile:        "/dev/null",
	})

	ctx := context.Background()
	_, err := client.RunCommand(ctx, "127.0.0.1", 2222, "root", "password", "echo test")

	assert.Error(t, err)
}

func TestClient_RunCommand_WithCustomPort(t *testing.T) {
	client := ssh.NewClient(ssh.Config{
		Timeout:               2 * time.Second,
		StrictHostKeyChecking: false,
		KnownHostsFile:        "",
	})

	ctx := context.Background()
	_, err := client.RunCommand(ctx, "127.0.0.1", 9999, "root", "password", "echo test")

	assert.Error(t, err)
}
