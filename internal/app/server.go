package app

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/auth"
	"TaskFlow/internal/handlers/comments"
	"TaskFlow/internal/handlers/middleware"
	"TaskFlow/internal/handlers/roles"
	"TaskFlow/internal/handlers/tasks"
	authUC "TaskFlow/internal/use-cases/auth"
	commentsUC "TaskFlow/internal/use-cases/comment"
	rolesUC "TaskFlow/internal/use-cases/roles"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(
	logger *slog.Logger,
	cfg *config.Config,
	authUc authUC.UseCases,
	tasksUc tasksUC.UseCases,
	commentsUC commentsUC.UseCases,
	checker middleware.PermissionChecker,
	rolesUc rolesUC.UseCases,
) *Server {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(SlogMiddleware(logger))
	router.Use(cors.New(cors.Config{
		AllowOrigins: splitAndTrimCORS(cfg.CORS.AllowedOrigins),
		AllowMethods: splitAndTrimCORS(cfg.CORS.AllowedMethods),
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	err := router.SetTrustedProxies(nil)
	if err != nil {
		logger.Error("set trusted proxies", slog.String("error", err.Error()))
	}

	auth.SetupRoutes(router, authUc, logger)
	tasks.SetupRoutes(router, cfg, tasksUc, checker, logger)
	comments.SetupRoutes(router, cfg, commentsUC, checker, logger)
	roles.SetupRoutes(router, cfg, rolesUc, checker, logger)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	appAddr := net.JoinHostPort(cfg.Listen.BindIP, cfg.Listen.Port)

	return &Server{
		httpServer: &http.Server{
			Addr:         appAddr,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Run() {
	s.logger.Info("starting server", slog.String("listen", s.httpServer.Addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("listen", slog.String("error", err.Error()))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	s.logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown", slog.String("error", err.Error()))
	}

	s.logger.Info("server exiting")
}

func SlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("client_ip", c.ClientIP()),
			slog.Int("body_size", c.Writer.Size()),
		)

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("request error", slog.String("error", err.Error()))
			}
		}
	}
}
