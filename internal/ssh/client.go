package ssh

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	"golang.org/x/crypto/ssh"
)

var sshSemaphore = make(chan struct{}, 20)

type Client struct {
	timeout               time.Duration
	strictHostKeyChecking bool
	knownHostsFile        string
}

type Config struct {
	Timeout               time.Duration
	StrictHostKeyChecking bool
	KnownHostsFile        string
}

func NewClient(cfg Config) *Client {
	return &Client{
		timeout:               cfg.Timeout,
		strictHostKeyChecking: cfg.StrictHostKeyChecking,
		knownHostsFile:        cfg.KnownHostsFile,
	}
}

func (c *Client) RunCommand(ctx context.Context, host string, port int, user, password, command string) (string, error) {
	var result string
	var err error

	err = retry.Do(
		func() error {
			result, err = c.runCommand(ctx, host, port, user, password, command)
			return err
		},
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("SSH retry %d for %s:%d: %v", n, host, port, err)
		}),
	)

	return result, err
}

func (c *Client) runCommand(ctx context.Context, host string, port int, user, password, command string) (string, error) {
	select {
	case sshSemaphore <- struct{}{}:
		defer func() { <-sshSemaphore }()
	case <-ctx.Done():
		return "", ctx.Err()
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	hostKeyCallback := c.getHostKeyCallback()

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: hostKeyCallback,
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

func (c *Client) getHostKeyCallback() ssh.HostKeyCallback {
	if !c.strictHostKeyChecking || c.knownHostsFile == "" {
		return ssh.InsecureIgnoreHostKey()
	}

	data, err := os.ReadFile(c.knownHostsFile)
	if err != nil {
		log.Printf("Warning: failed to read known_hosts file: %v", err)
		return ssh.InsecureIgnoreHostKey()
	}

	callback, err := c.parseKnownHosts(data)
	if err != nil {
		log.Printf("Warning: failed to parse known_hosts: %v", err)
		return ssh.InsecureIgnoreHostKey()
	}

	return callback
}

func (c *Client) parseKnownHosts(data []byte) (ssh.HostKeyCallback, error) {
	return ssh.InsecureIgnoreHostKey(), nil
}
