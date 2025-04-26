package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bariscan97/clean-rest-architecture/app/routes"
	post_handler "github.com/bariscan97/clean-rest-architecture/internal/handler/post"
	user_handler "github.com/bariscan97/clean-rest-architecture/internal/handler/user"
	post_repo "github.com/bariscan97/clean-rest-architecture/internal/repository/post"
	user_repo "github.com/bariscan97/clean-rest-architecture/internal/repository/user"
	"github.com/bariscan97/clean-rest-architecture/pkg/config"
	"github.com/bariscan97/clean-rest-architecture/pkg/database"
	"github.com/ianschenck/envflag"
	"go.uber.org/zap"
)

const minSecretKeySize = 32

var (
	addr string
	secretKey *string
)


func gracefulShutdown(server *http.Server) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	zap.L().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}

func main() {
	cfg, _ := config.LoadConfig("*")
	db := database.NewConnection(cfg)

	addr = strconv.Itoa(cfg.App.Port)
	secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for JWT signing")
	
	if len(*secretKey) < minSecretKeySize {
		zap.L().Fatal("SECRET_KEY must be at least X characters",
            zap.Int("minSecretKeySize", minSecretKeySize),
            zap.Int("gotLength", len(*secretKey)),
        )
	}

	userRepo := user_repo.NewUserRepository(db)
	postRepo := post_repo.NewUserRepository(db)

	userHandler := user_handler.NewUserHandler(userRepo, *secretKey)
	postHandler := post_handler.NewPostHandler(postRepo)

	r := routes.NewRouter(
		*userHandler,
		*postHandler,
	)
	r.RegisterRoutes()

	srv := &http.Server{
		Addr:         addr,
		Handler:      r.Mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	zap.L().Info("Server started on port", zap.String("port", addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Error while starting server", zap.Error(err))
		}
	}()

	gracefulShutdown(srv)
}
