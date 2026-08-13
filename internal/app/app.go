package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/config"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/repository/postgres"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/server"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/service"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/ssh"
)

type App struct {
	Config *config.Config
	Logger *slog.Logger
	Server *server.Server
}

func New() *App {
	cfg, logger := initConfigAndLogger()
	db := initDatabase(cfg, logger)

	repos := initRepositories(db)
	services := initServices(repos, cfg)
	srv := server.New(services, cfg, logger)

	return &App{
		Config: cfg,
		Logger: logger,
		Server: srv,
	}
}

func (a *App) Run() {
	a.Server.Start(a.Logger)
	a.waitForShutdown()
	a.shutdown()
}

func (a *App) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (a *App) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.Server.Shutdown(ctx, a.Logger)
}

func initConfigAndLogger() (*config.Config, *slog.Logger) {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.Log)
	logger.Info("config loaded",
		"create_port", cfg.Server.CreatePort,
		"callback_port", cfg.Server.CallbackPort,
		"db_host", cfg.Database.Host,
		"db_port", cfg.Database.Port,
		"db_name", cfg.Database.DBName,
	)

	return cfg, logger
}

func initDatabase(cfg *config.Config, logger *slog.Logger) *postgres.DB {
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connected successfully")
	return db
}

type Repositories struct {
	ScriptRepo *postgres.ScriptRepository
	EventRepo  *postgres.EventRepository
}

func initRepositories(db *postgres.DB) *Repositories {
	return &Repositories{
		ScriptRepo: postgres.NewScriptRepository(db),
		EventRepo:  postgres.NewEventRepository(db),
	}
}

func initServices(repos *Repositories, cfg *config.Config) *service.Services {
	sshClient := ssh.NewClient(cfg.SSH.Timeout)

	return &service.Services{
		ScriptService: service.NewScriptService(repos.ScriptRepo, sshClient),
		EventService:  service.NewEventService(repos.EventRepo, repos.ScriptRepo),
	}
}

func setupLogger(cfg config.LogConfig) *slog.Logger {
	var handler slog.Handler

	level := slog.LevelInfo
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
