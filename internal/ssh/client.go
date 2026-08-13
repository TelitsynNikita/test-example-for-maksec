package ssh

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	config  *ssh.ClientConfig
	timeout time.Duration
}

type Config struct {
	User     string
	Password string
	Timeout  time.Duration
}

func NewClient(cfg Config) *Client {
	return &Client{
		config: &ssh.ClientConfig{
			User: cfg.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(cfg.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         cfg.Timeout,
		},
		timeout: cfg.Timeout,
	}
}

func (c *Client) RunCommand(ctx context.Context, host string, port int, command string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := ssh.Dial("tcp", addr, c.config)
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
