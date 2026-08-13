package ssh

import "context"

type SSHClient interface {
	RunCommand(ctx context.Context, host string, port int, command string) (string, error)
}
