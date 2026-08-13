package ssh

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// Ограничиваем количество параллельных SSH подключений
var sshSemaphore = make(chan struct{}, 20)

type Client struct {
	timeout time.Duration
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		timeout: timeout,
	}
}

func (c *Client) RunCommand(ctx context.Context, host string, port int, user, password, command string) (string, error) {
	// Bulkhead: ограничиваем количество параллельных SSH подключений
	select {
	case sshSemaphore <- struct{}{}:
		defer func() { <-sshSemaphore }()
	case <-ctx.Done():
		return "", ctx.Err()
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}
