package api

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/interfaces/repositories"
	"go-learning/clean-architecture-mongo/usecases"
	"net/http"
	"time"
)

type UserController struct {
	UserInterceptor usecases.UserInterceptor
	Logger          usecases.Logger
}

func NewUserController(db database.Database, timeout time.Duration, logger usecases.Logger) *UserController {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)

	return &UserController{
		UserInterceptor: usecases.UserInterceptor{
			UserRepository: ur,
			ContextTimeout: timeout,
		},
		Logger: logger,
	}
}

func (uc *UserController) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"success": "true"})
	return
}
