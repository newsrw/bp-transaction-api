package server

import (
	"bp-transaction-api/configs"
	transactionHandler "bp-transaction-api/transaction/delivery"
	transactionUsecase "bp-transaction-api/transaction/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server server
type Server struct {
	port   string
	ctx    context.Context
	logger *zap.Logger
	config *configs.Config
}

// New new server
func New(config *configs.Config) (*Server, error) {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	restPort := fmt.Sprintf(":%d", config.Client.Port)

	return &Server{
		port:   restPort,
		ctx:    context.Background(),
		logger: logger,
		config: config,
	}, nil
}

// Start start server
func (s *Server) Start() {
	s.logger.Info("Start RESTful", zap.String("PORT", s.port))

	r := gin.Default()

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Healthcheck Successful",
		})
	})

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	tu := transactionUsecase.NewTransactionUsecase(timeoutContext, s.config, s.logger)
	transactionHandler.NewTransactionHandler(r, s.config, tu)

	srv := &http.Server{
		Addr:    s.port,
		Handler: r,
	}

	s.gracefullyShutdown(srv)
}

func (s *Server) gracefullyShutdown(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force shutdown")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exit")
}
