package http

import (
	"github.com/ReilEgor/CleanArchitectureGolang/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ginServer struct {
	router *gin.Engine
	uc     *usecase.TaskUsecase
}

func NewGinServer(uc *usecase.TaskUsecase) *ginServer {
	s := &ginServer{
		router: gin.Default(),
		uc:     uc,
	}
	s.mapRoutes()
	return s
}

func (s *ginServer) Run(port string) error {
	return s.router.Run(port)
}

func (s *ginServer) mapRoutes() {
	h := NewTaskHandler(s.uc)

	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/tasks", h.GetTasks)
		v1.POST("/tasks", h.CreateTask)
	}
}

func (s *ginServer) GetRouter() *gin.Engine {
	return s.router
}
