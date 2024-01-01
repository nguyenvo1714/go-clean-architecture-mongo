package auth

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/interfaces/repositories"
	"go-learning/clean-architecture-mongo/usecases"
	authUseCases "go-learning/clean-architecture-mongo/usecases/auth"
	"net/http"
	"time"
)

type ProfileController struct {
	ProfileInterceptor authUseCases.ProfileInterceptor
	Logger             usecases.Logger
	Config             *infrastructure.Config
}

func NewProfileController(config *infrastructure.Config, db database.Database, timeout time.Duration, logger usecases.Logger) *ProfileController {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)

	return &ProfileController{
		ProfileInterceptor: authUseCases.ProfileInterceptor{
			UserRepository: ur,
			ContextTimeout: timeout,
		},
		Logger: logger,
		Config: config,
	}
}

func (pc *ProfileController) Show(c *gin.Context) {
	userID := c.GetString("x-user-id")
	profile, err := pc.ProfileInterceptor.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, profile)
}
