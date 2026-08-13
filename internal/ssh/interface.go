package ssh

import "context"

type SSHClient interface {
	RunCommand(ctx context.Context, host string, port int, user, password, command string) (string, error)
}
