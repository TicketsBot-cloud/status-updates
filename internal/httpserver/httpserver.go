package httpserver

import (
	"github.com/TicketsBot-cloud/status-updates/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server represents the HTTP server with configuration and logging.
type Server struct {
	logger *zap.Logger   // logger is used for structured logging
	config config.Config // config holds the server configuration
}

// NewServer creates a new Server instance with the provided logger and configuration.
func NewServer(logger *zap.Logger, conf config.Config) *Server {
	return &Server{
		logger: logger,
		config: conf,
	}
}

// Start launches the HTTP server and sets up the routes.
// It returns an error if the server fails to start.
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.String("address", s.config.ServerAddr))

	router := gin.New()

	router.POST("/interactions", s.AuthMiddleware, s.HandleInteraction)

	return router.Run(s.config.ServerAddr)
}
