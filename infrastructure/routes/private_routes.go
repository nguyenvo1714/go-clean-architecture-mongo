package routes

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/interfaces/api"
	"go-learning/clean-architecture-mongo/interfaces/api/auth"
	"go-learning/clean-architecture-mongo/usecases"
	"time"
)

func DispatchPrivateRoutes(group *gin.RouterGroup, logger usecases.Logger, db database.Database, timeout time.Duration, config *infrastructure.Config) {
	taskController := api.NewTaskController(db, timeout, logger)
	profileController := auth.NewProfileController(config, db, timeout, logger)

	group.GET("/profile", profileController.Show)
	group.GET("/tasks", taskController.Index)
	group.POST("tasks/create", taskController.Store)
}
