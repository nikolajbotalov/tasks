package main

import (
	"TaskFlow/internal/config"
	router "TaskFlow/internal/delivery/http"
	authRepo "TaskFlow/internal/repository/auth"
	commentsRepo "TaskFlow/internal/repository/comment"
	permissionsRepo "TaskFlow/internal/repository/permissions"
	"TaskFlow/internal/repository/postgres"
	tasksRepo "TaskFlow/internal/repository/tasks"
	authUC "TaskFlow/internal/use-cases/auth"
	commentsUC "TaskFlow/internal/use-cases/comment"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBName, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBSSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ping db: %v", err)
	}

	log.Println("Successfully connected to db")

	authRepository := authRepo.NewAuthRepository(db)
	authUseCases := authUC.NewAuthUseCases(authRepository, cfg)
	taskRepository := tasksRepo.NewTasksRepository(db)
	tasksUseCases := tasksUC.NewTasksUseCases(taskRepository)
	commentsRepository := commentsRepo.NewCommentRepository(db)
	commentsUseCases := commentsUC.NewCommentUseCases(commentsRepository)

	permissionsRepository := permissionsRepo.NewPermissionsRepository(db)
	rolePermissions, err := permissionsRepository.GetAll()
	if err != nil {
		log.Printf("load role permissions: %v\n", err)
	}

	r := router.AppRouters(authUseCases, cfg, tasksUseCases, commentsUseCases, rolePermissions)

	ctx := context.Background()

	if err = postgres.EnsureSchema(ctx, db, "tasks", "users", "task_comments", "roles", "role_permissions"); err != nil {
		log.Printf("schema check: %v\n", err)
	}

	appAddr := net.JoinHostPort(cfg.Listen.BindIP, cfg.Listen.Port)
	err = r.Run(appAddr)
	if err != nil {
		log.Fatalf("start server: %v", err)
	}
}
