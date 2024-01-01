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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type SignupController struct {
	SignupInterceptor authUseCases.SignupInterceptor
	Logger            usecases.Logger
	Config            *infrastructure.Config
}

func NewSignupController(config *infrastructure.Config, db database.Database, timeout time.Duration, logger usecases.Logger) *SignupController {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)

	return &SignupController{
		SignupInterceptor: authUseCases.SignupInterceptor{
			UserRepository: ur,
			ContextTimeout: timeout,
		},
		Logger: logger,
		Config: config,
	}
}

func (sc *SignupController) Signup(c *gin.Context) {
	var data authDomain.Signup
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ErrorResponse{Message: err.Error()})

		return
	}

	_, err := sc.SignupInterceptor.GetUserByEmail(c, data.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "User already exists."})

		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	data.Password = string(encryptedPassword)
	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}

	err = sc.SignupInterceptor.Store(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, err := sc.SignupInterceptor.CreateAccessToken(user, sc.Config.AccessTokenSecret, sc.Config.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := sc.SignupInterceptor.CreateRefreshToken(user, sc.Config.RefreshTokenSecret, sc.Config.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, authDomain.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
