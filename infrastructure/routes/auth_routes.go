package routes

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/interfaces/api/auth"
	"go-learning/clean-architecture-mongo/usecases"
	"time"
)

func DispatchAuthRoutes(group *gin.RouterGroup, logger usecases.Logger, db database.Database, timeout time.Duration, config *infrastructure.Config) {
	loginController := auth.NewLoginController(config, db, timeout, logger)
	signupController := auth.NewSignupController(config, db, timeout, logger)
	refreshTokenController := auth.NewRefreshTokenController(config, db, timeout, logger)

	group.POST("/login", loginController.Authenticate)
	group.POST("/signup", signupController.Signup)
	group.POST("/token-refresh", refreshTokenController.RefreshToken)
}
