package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/handler"
)

type Router struct {
	handlers *handler.Handler
}

func NewRouter(handlers *handler.Handler) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// API v1 group
	v1 := router.Group("/api")
	{
		// Register all route handlers
		r.handlers.Schedule.RegisterRoutes(v1)
		// Add other route handlers here as needed
	}

	return router
}
