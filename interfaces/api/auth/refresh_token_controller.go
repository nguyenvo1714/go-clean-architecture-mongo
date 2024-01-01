package auth

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	authDomain "go-learning/clean-architecture-mongo/domain/auth"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/interfaces/repositories"
	"go-learning/clean-architecture-mongo/usecases"
	authUseCases "go-learning/clean-architecture-mongo/usecases/auth"
	"go-learning/clean-architecture-mongo/utils"
	"net/http"
	"time"
)

type RefreshTokenController struct {
	RefreshTokenInterceptor authUseCases.RefreshTokenInterceptor
	Logger                  usecases.Logger
	Config                  *infrastructure.Config
}

func NewRefreshTokenController(config *infrastructure.Config, db database.Database, timeout time.Duration, logger usecases.Logger) *RefreshTokenController {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)

	return &RefreshTokenController{
		RefreshTokenInterceptor: authUseCases.RefreshTokenInterceptor{
			UserRepository: ur,
			ContextTimeout: timeout,
		},
		Logger: logger,
		Config: config,
	}
}

func (rc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var requestToken authDomain.RefreshTokenRequest

	if err := c.ShouldBind(&requestToken); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID, err := utils.ExtractIDFromToken(requestToken.RefreshToken, rc.Config.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := rc.RefreshTokenInterceptor.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, err := utils.CreateAccessToken(&user, rc.Config.AccessTokenSecret, rc.Config.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := utils.CreateRefreshToken(&user, rc.Config.RefreshTokenSecret, rc.Config.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, authDomain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
