package app

import (
	"TaskFlow/internal/config"
	authRepo "TaskFlow/internal/repository/auth"
	commentsRepo "TaskFlow/internal/repository/comment"
	permissionsRepo "TaskFlow/internal/repository/permissions"
	"TaskFlow/internal/repository/postgres"
	rolesRepo "TaskFlow/internal/repository/roles"
	tasksRepo "TaskFlow/internal/repository/tasks"
	authUC "TaskFlow/internal/use-cases/auth"
	commentsUC "TaskFlow/internal/use-cases/comment"
	rolesUC "TaskFlow/internal/use-cases/roles"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"TaskFlow/migrations"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type App struct {
	Logger *slog.Logger
	Cfg    *config.Config
	DB     *sql.DB
	Server *Server
}

func NewApp() (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("read config", slog.String("error", err.Error()))
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBName, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBSSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("connect db", slog.String("error", err.Error()))
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Minute)

	err = db.Ping()
	if err != nil {
		logger.Error("ping db", slog.String("error", err.Error()))
		return nil, err
	}

	if err = postgres.RunMigrations(context.Background(), db, migrations.FS); err != nil {
		logger.Error("migrate db", slog.String("error", err.Error()))
		return nil, err
	}

	authRepository := authRepo.NewAuthRepository(db)
	authUseCases := authUC.NewAuthUseCases(authRepository, cfg)
	taskRepository := tasksRepo.NewTasksRepository(db)
	tasksUseCases := tasksUC.NewTasksUseCases(taskRepository)
	commentsRepository := commentsRepo.NewCommentRepository(db)
	commentsUseCases := commentsUC.NewCommentUseCases(commentsRepository)
	rolesRepository := rolesRepo.NewRolesRepository(db)
	rolesUseCases := rolesUC.NewRolesUseCases(rolesRepository)

	permissionsRepository := permissionsRepo.NewPermissionsRepository(db)
	rolePermissions, err := permissionsRepository.GetAll()
	if err != nil {
		logger.Error("load role permissions", slog.String("error", err.Error()))
		return nil, err
	}

	server := NewServer(logger, cfg, authUseCases, tasksUseCases, commentsUseCases, rolePermissions, rolesUseCases)

	return &App{
		Logger: logger,
		Cfg:    cfg,
		DB:     db,
		Server: server,
	}, nil
}

func (a *App) Close() {
	a.Logger.Info("closing app")
	if err := a.DB.Close(); err != nil {
		a.Logger.Error("close db", slog.String("error", err.Error()))
	}
}
