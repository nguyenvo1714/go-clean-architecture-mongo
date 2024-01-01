package api

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/interfaces/repositories"
	"go-learning/clean-architecture-mongo/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type TaskController struct {
	TaskInterceptor usecases.TaskInterceptor
	Logger          usecases.Logger
}

func NewTaskController(db database.Database, timeout time.Duration, logger usecases.Logger) *TaskController {
	tr := repositories.NewTaskRepository(db, domain.CollectionTask)

	return &TaskController{
		TaskInterceptor: usecases.TaskInterceptor{
			TaskRepository: tr,
			ContextTimeout: timeout,
		},
		Logger: logger,
	}
}

func (tc *TaskController) Index(c *gin.Context) {
	userID := c.GetString("x-user-id")
	tasks, err := tc.TaskInterceptor.FetchByUserID(c, userID)
	if err != nil {
		tc.Logger.LogError("%s", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) Store(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		tc.Logger.LogError("%s", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	userID := c.GetString("x-user-id")
	task.ID = primitive.NewObjectID()

	task.UserID, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		tc.Logger.LogError("%s", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	err = tc.TaskInterceptor.Create(c, &task)
	if err != nil {
		tc.Logger.LogError("%s", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully",
	})
}
