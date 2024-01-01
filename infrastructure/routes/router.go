package routes

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/infrastructure/middleware"
	"go-learning/clean-architecture-mongo/usecases"
	"time"
)

func Dispatch(router *gin.Engine, config *infrastructure.Config, db database.Database, logger usecases.Logger, timeout time.Duration) {
	authRouter := router.Group("/v1")
	DispatchAuthRoutes(authRouter, logger, db, timeout, config)

	protectedRouter := router.Group("/v1")
	protectedRouter.Use(middleware.JwtAuthMiddleware(config.AccessTokenSecret))
	DispatchPrivateRoutes(protectedRouter, logger, db, timeout, config)
}
