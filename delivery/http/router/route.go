package router

import (
	"auth-service/domain/service"

	"auth-service/application/usecase"
	"auth-service/config"
	"auth-service/delivery/http/handler"
	"auth-service/delivery/http/middlewares"
	"auth-service/delivery/http/request"
	"auth-service/infrastructure/keycloack"
	"auth-service/infrastructure/persistence/mysql"
	pkgAppHttp "auth-service/pkg/app_http"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewRouter(
	cfg config.Config,
	txManager *mysql.TransactionManager,
	db *gorm.DB,
	validator *validator.Validate,
	httpClient *pkgAppHttp.AppHttp,
	ut ut.Translator,
) *mux.Router {

	router := mux.NewRouter()

	// global middleware
	router.Use(
		middlewares.CORS,
		middlewares.JSON,
	)

	// repository
	userRepo := mysql.NewUserRepository(db)

	// infrastructure service
	keycloakService := keycloack.NewKeycloakService(
		httpClient,
		cfg.Keycloak,
	)

	keycloakAdapter := keycloack.NewKeycloakAdapter(
		keycloakService,
	)

	// service
	userService := service.NewUserService(
		userRepo,
		txManager,
		cfg,
		keycloakAdapter,
	)

	// usecase
	createUserUseCase := usecase.NewCreateUserUseCase(userService)
	listUserUseCase := usecase.NewListUserUseCase(userRepo)
	loginUserUseCase := usecase.NewLoginUserUseCase(userService, userRepo)
	refreshTokenUseCase := usecase.NewRefreshTokenUserUserUseCase(userService)
	getUserUseCase := usecase.NewDetailUserUseCase(userRepo)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userService)
	deleteUserUseCase := usecase.NewDeleteUserUseCase(userService)

	// handler
	createUserHandler := handler.NewCreateUserHandler(createUserUseCase)
	listUserHandler := handler.NewListUserHandler(listUserUseCase)
	loginUserHandler := handler.NewLoginUserHandler(loginUserUseCase)
	refreshTokenHandler := handler.NewRefreshTokenUserHandler(refreshTokenUseCase)
	getUserHandler := handler.NewGetUserHandler(getUserUseCase)
	updateUserHandler := handler.NewUpdateUserHandler(updateUserUseCase)
	deleteUserHandler := handler.NewDeleteUserHandler(deleteUserUseCase)

	// health check
	router.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}).Methods(http.MethodGet)

	// api v1
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// public routes
	users := v1.PathPrefix("/users").Subrouter()

	users.Handle(
		"/login",
		middlewares.ValidateRequestBody[request.LoginUserRequest](validator, ut)(
			http.HandlerFunc(loginUserHandler.Execute),
		),
	).Methods(http.MethodPost)

	users.Handle(
		"/refresh-token",
		middlewares.ValidateRequestBody[request.RefreshTokenRequest](validator, ut)(
			http.HandlerFunc(refreshTokenHandler.Execute),
		),
	).Methods(http.MethodPost)

	// protected routes
	privateUsers := v1.PathPrefix("/users").Subrouter()
	privateUsers.Use(middlewares.Auth)

	privateUsers.Handle(
		"",
		middlewares.ValidateRequestBody[request.CreateUserRequest](validator, ut)(
			http.HandlerFunc(createUserHandler.Execute),
		),
	).Methods(http.MethodPost)

	privateUsers.HandleFunc(
		"",
		listUserHandler.Execute,
	).Methods(http.MethodGet)

	privateUsers.HandleFunc(
		"/{id}",
		getUserHandler.Execute,
	).Methods(http.MethodGet)

	privateUsers.Handle(
		"/{id}",
		middlewares.ValidateRequestBody[request.UpdateUserRequest](validator, ut)(
			http.HandlerFunc(updateUserHandler.Execute),
		),
	).Methods(http.MethodPut)

	privateUsers.HandleFunc(
		"/{id}",
		deleteUserHandler.Execute,
	).Methods(http.MethodDelete)

	return router
}

func RunHTTPServer(
	router *mux.Router,
	port string,
	logger *logrus.Logger,
) error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           handlers.CompressHandler(router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Channel to capture server errors
	serverErrCh := make(chan error, 1)

	// Start server in goroutine
	go func() {
		logger.Infof("HTTP server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			serverErrCh <- err
		}
	}()

	select {
	case sig := <-quit:
		logger.Infof("Received OS signal: %s, shutting down...", sig)
	case err := <-serverErrCh:
		return fmt.Errorf("HTTP server error: %w", err)
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logger.Info("Server exited gracefully")
	return nil
}
