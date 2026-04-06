package router

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/handler"
	"github.com/Kyouheip/MathOvercome_serverless/internal/middleware"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func New(client *dynamodb.Client) *gin.Engine {
	repo := repository.NewRepository(client)
	testSessSvc := service.NewTestSessionService(repo)
	mypageSvc := service.NewMypageService(repo)
	sessionHandler := handler.NewSessionHandler(testSessSvc, mypageSvc)

	r := gin.Default()

	if os.Getenv("APP_ENV") == "local" {
		allowOrigin := os.Getenv("ALLOW_ORIGIN")
		r.Use(middleware.LocalAuthMiddleware())
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{allowOrigin},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}



	sess := r.Group("/session")
	{
		sess.POST("/test", sessionHandler.CreateTestSess)
		sess.GET("/current/problems/:idx", sessionHandler.ViewOneProblem)
		sess.POST("/current/problems/:idx/answer", sessionHandler.SubmitAnswer)
		sess.GET("/mypage", sessionHandler.GetMypage)
	}

	return r
}
