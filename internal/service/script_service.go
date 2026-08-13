package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/repository"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/ssh"
	"github.com/google/uuid"
)

type ScriptService struct {
	scriptRepo repository.ScriptRepository
	sshClient  ssh.SSHClient
	templates  map[string]domain.ScriptTemplate
}

func NewScriptService(scriptRepo repository.ScriptRepository, sshClient ssh.SSHClient) *ScriptService {
	return &ScriptService{
		scriptRepo: scriptRepo,
		sshClient:  sshClient,
		templates:  domain.GetTemplates(),
	}
}

func (s *ScriptService) CreateScript(ctx context.Context, req domain.CreateScriptRequest) (*domain.Script, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	templateContent, err := s.getTemplateContent(req.Template)
	if err != nil {
		return nil, err
	}

	scriptID := uuid.New()
	scriptPath := fmt.Sprintf("/opt/script-monitor/scripts/%s.sh", scriptID.String())

	exists, err := s.scriptRepo.Exists(ctx, scriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check script existence: %w", err)
	}
	if exists {
		return nil, domain.ErrScriptAlreadyExists
	}

	sshPort := 22
	createDirCmd := "mkdir -p /opt/script-monitor/scripts"
	if _, err := s.sshClient.RunCommand(ctx, req.Host, sshPort, createDirCmd); err != nil {
		return nil, fmt.Errorf("failed to create remote directory: %w", err)
	}

	escapedContent := strings.ReplaceAll(templateContent, "'", "'\\''")
	uploadCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", scriptPath, escapedContent)
	if _, err := s.sshClient.RunCommand(ctx, req.Host, sshPort, uploadCmd); err != nil {
		return nil, fmt.Errorf("failed to upload script: %w", err)
	}

	chmodCmd := fmt.Sprintf("chmod +x %s", scriptPath)
	if _, err := s.sshClient.RunCommand(ctx, req.Host, sshPort, chmodCmd); err != nil {
		s.cleanupRemoteFile(ctx, req.Host, scriptPath)
		return nil, fmt.Errorf("failed to make script executable: %w", err)
	}

	script := domain.NewScript(req.Host, req.User, req.Template, scriptPath)
	if err := s.scriptRepo.Create(ctx, script); err != nil {
		s.cleanupRemoteFile(ctx, req.Host, scriptPath)
		return nil, fmt.Errorf("failed to save script to database: %w", err)
	}

	return script, nil
}

func (s *ScriptService) validateRequest(req domain.CreateScriptRequest) error {
	if req.Host == "" {
		return domain.ErrInvalidHost
	}
	if req.User == "" {
		return domain.ErrInvalidUser
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

func (s *ScriptService) getTemplateContent(name string) (string, error) {
	tmpl, exists := s.templates[name]
	if !exists {
		return "", domain.ErrTemplateNotFound
	}
	return tmpl.Content, nil
}

func (s *ScriptService) cleanupRemoteFile(ctx context.Context, host, path string) {
	cmd := fmt.Sprintf("rm -f %s", path)
	s.sshClient.RunCommand(ctx, host, 22, cmd)
}
