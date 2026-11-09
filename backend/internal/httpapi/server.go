package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	db     *mongo.Database
	logger *slog.Logger
}

func New(client *mongo.Client, database string, logger *slog.Logger) http.Handler {
	server := &Server{db: client.Database(database), logger: logger}
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/healthz", server.health)
	router.GET("/readyz", server.ready)
	return router
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := s.db.Client().Ping(ctx, nil); err != nil {
		s.logger.Warn("readiness check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
