package app

import (
	"TaskFlow/internal/config"
	authRepo "TaskFlow/internal/repository/auth"
	commentsRepo "TaskFlow/internal/repository/comment"
	permissionsRepo "TaskFlow/internal/repository/permissions"
	"TaskFlow/internal/repository/postgres"
	profilesRepo "TaskFlow/internal/repository/profiles"
	rolesRepo "TaskFlow/internal/repository/roles"
	tasksRepo "TaskFlow/internal/repository/tasks"
	authUC "TaskFlow/internal/use-cases/auth"
	commentsUC "TaskFlow/internal/use-cases/comment"
	profilesUC "TaskFlow/internal/use-cases/profiles"
	rolesUC "TaskFlow/internal/use-cases/roles"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"database/sql"
	"log/slog"
	"os"
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

	db, err := postgres.ConnectDB(cfg)
	if err != nil {
		logger.Error("connect database", slog.String("error", err.Error()))
		return nil, err
	}

	authRepository := authRepo.NewAuthRepository(db)
	authUseCases := authUC.NewAuthUseCases(authRepository, cfg)
	profileRepository := profilesRepo.NewProfilesRepository(db)
	profilesUseCases := profilesUC.NewProfilesUseCases(profileRepository)
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

	server := NewServer(
		logger,
		cfg,
		authUseCases,
		tasksUseCases,
		commentsUseCases,
		rolePermissions,
		rolesUseCases,
		profilesUseCases,
	)

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
