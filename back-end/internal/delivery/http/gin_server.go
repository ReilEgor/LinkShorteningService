package http

import (
	"log/slog"

	"github.com/ReilEgor/LinkShorteningService/internal/delivery/http/handler"
	"github.com/ReilEgor/LinkShorteningService/internal/delivery/http/middleware"
	"github.com/ReilEgor/LinkShorteningService/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ginServer struct {
	router *gin.Engine
	uc     *usecase.LinkUsecase
	logger *slog.Logger
}

func NewGinServer(uc *usecase.LinkUsecase, logger *slog.Logger) *ginServer {
	s := &ginServer{
		router: gin.New(),
		uc:     uc,
		logger: logger,
	}
	s.router.Use(gin.Recovery())

	s.router.Use(middleware.RequestIDMiddleware())
	s.mapRoutes()
	return s
}

func (s *ginServer) Run(port string) error {
	return s.router.Run(port)
}

func (s *ginServer) mapRoutes() {
	h := handler.NewLinkHandler(s.uc, s.logger)
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/:shortURL", h.GetLink)
		v1.POST("/longURL", h.AddLink)
	}
}

func (s *ginServer) GetRouter() *gin.Engine {
	return s.router
}
