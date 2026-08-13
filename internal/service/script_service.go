package service

import (
	"context"
	"fmt"

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

	sshPort := req.Port
	if sshPort == 0 {
		sshPort = 22
	}

	// 1. Создаем директорию
	createDirCmd := "mkdir -p /opt/script-monitor/scripts"
	if _, err := s.sshClient.RunCommand(ctx, req.Host, sshPort, req.User, req.Password, createDirCmd); err != nil {
		return nil, fmt.Errorf("failed to create remote directory: %w", err)
	}

	// 2. Загружаем скрипт
	uploadCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", scriptPath, templateContent)
	if _, err := s.sshClient.RunCommand(ctx, req.Host, sshPort, req.User, req.Password, uploadCmd); err != nil {
		s.cleanupRemoteFile(ctx, req.Host, scriptPath, req.User, req.Password, sshPort)
		return nil, fmt.Errorf("failed to upload script: %w", err)
	}

	// 3. Делаем скрипт исполняемым
	chmodCmd := fmt.Sprintf("chmod +x %s", scriptPath)
	if _, err := s.sshClient.RunCommand(ctx, req.Host, sshPort, req.User, req.Password, chmodCmd); err != nil {
		s.cleanupRemoteFile(ctx, req.Host, scriptPath, req.User, req.Password, sshPort)
		return nil, fmt.Errorf("failed to make script executable: %w", err)
	}

	// 4. Устанавливаем агент мониторинга
	if err := s.installAgent(ctx, req.Host, sshPort, req.User, req.Password); err != nil {
		s.cleanupRemoteFile(ctx, req.Host, scriptPath, req.User, req.Password, sshPort)
		return nil, fmt.Errorf("failed to install agent: %w", err)
	}

	// 5. Сохраняем в БД
	script := domain.NewScript(req.Host, req.User, req.Template, scriptPath)
	if err := s.scriptRepo.Create(ctx, script); err != nil {
		s.cleanupRemoteFile(ctx, req.Host, scriptPath, req.User, req.Password, sshPort)
		return nil, fmt.Errorf("failed to save script to database: %w", err)
	}

	return script, nil
}

func (s *ScriptService) installAgent(ctx context.Context, host string, port int, user, password string) error {
	// Проверяем, установлен ли агент
	checkCmd := "test -f /opt/script-monitor/agent/agent.sh"
	if _, err := s.sshClient.RunCommand(ctx, host, port, user, password, checkCmd); err == nil {
		// Агент уже установлен
		return nil
	}

	// Создаем директорию для агента
	mkdirCmd := "mkdir -p /opt/script-monitor/agent"
	if _, err := s.sshClient.RunCommand(ctx, host, port, user, password, mkdirCmd); err != nil {
		return fmt.Errorf("failed to create agent directory: %w", err)
	}

	// Получаем содержимое агента
	agentContent := domain.GetAgentScript()

	// Загружаем агент
	uploadCmd := fmt.Sprintf("cat > /opt/script-monitor/agent/agent.sh << 'EOF'\n%s\nEOF", agentContent)
	if _, err := s.sshClient.RunCommand(ctx, host, port, user, password, uploadCmd); err != nil {
		return fmt.Errorf("failed to upload agent: %w", err)
	}

	// Делаем агент исполняемым
	chmodCmd := "chmod +x /opt/script-monitor/agent/agent.sh"
	if _, err := s.sshClient.RunCommand(ctx, host, port, user, password, chmodCmd); err != nil {
		return fmt.Errorf("failed to make agent executable: %w", err)
	}

	// Запускаем агент в фоне
	startCmd := "nohup /opt/script-monitor/agent/agent.sh > /opt/script-monitor/agent/agent.log 2>&1 &"
	if _, err := s.sshClient.RunCommand(ctx, host, port, user, password, startCmd); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	return nil
}

func (s *ScriptService) validateRequest(req domain.CreateScriptRequest) error {
	if req.Host == "" {
		return domain.ErrInvalidHost
	}
	if len(req.Host) > 255 {
		return fmt.Errorf("host must be less than 255 characters")
	}
	if req.User == "" {
		return domain.ErrInvalidUser
	}
	if len(req.User) > 100 {
		return fmt.Errorf("user must be less than 100 characters")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) > 128 {
		return fmt.Errorf("password must be less than 128 characters")
	}
	if req.Template == "" {
		return fmt.Errorf("template is required")
	}
	if len(req.Template) > 50 {
		return fmt.Errorf("template must be less than 50 characters")
	}
	if req.Port < 1 || req.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
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

func (s *ScriptService) cleanupRemoteFile(ctx context.Context, host, path, user, password string, port int) {
	cmd := fmt.Sprintf("rm -f %s", path)
	s.sshClient.RunCommand(ctx, host, port, user, password, cmd)
}
