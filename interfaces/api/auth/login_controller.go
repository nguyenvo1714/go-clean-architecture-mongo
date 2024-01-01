package auth

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	authDomain "go-learning/clean-architecture-mongo/domain/auth"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/interfaces/repositories"
	"go-learning/clean-architecture-mongo/usecases"
	authUsecases "go-learning/clean-architecture-mongo/usecases/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginController struct {
	LoginInterceptor authUsecases.LoginInterceptor
	Logger           usecases.Logger
	Config           *infrastructure.Config
}

func NewLoginController(config *infrastructure.Config, db database.Database, timeout time.Duration, logger usecases.Logger) *LoginController {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)

	return &LoginController{
		LoginInterceptor: authUsecases.LoginInterceptor{
			UserRepository: ur,
			ContextTimeout: timeout,
		},
		Logger: logger,
		Config: config,
	}
}

func (lc *LoginController) Authenticate(c *gin.Context) {
	var credentials authDomain.LoginRequest

	err := c.ShouldBind(&credentials)
	if err != nil {
		lc.Logger.LogError("%s", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	user, err := lc.LoginInterceptor.GetUserByEmail(c, credentials.Email)
	if err != nil {
		lc.Logger.LogError("%s", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		lc.Logger.LogError("%s", err)
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Message: "Invalid credentials",
		})

		return
	}

	accessToken, err := lc.LoginInterceptor.CreateAccessToken(
		&user,
		lc.Config.AccessTokenSecret,
		lc.Config.AccessTokenExpiryHour,
	)
	if err != nil {
		lc.Logger.LogError("%s", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	refreshToken, err := lc.LoginInterceptor.CreateRefreshToken(
		&user,
		lc.Config.RefreshTokenSecret,
		lc.Config.RefreshTokenExpiryHour,
	)
	if err != nil {
		lc.Logger.LogError("%s", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, authDomain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
