package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/handler"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"

	"gorm.io/gorm"
)

func New(db *gorm.DB, sessionSecret string) *gin.Engine {
	repo := repository.NewRepository(db)
	loginSvc := service.NewLoginService(repo)
	testSessSvc := service.NewTestSessionService(repo)
	mypageSvc := service.NewMypageService(repo)
	authHandler := handler.NewAuthHandler(loginSvc)
	sessionHandler := handler.NewSessionHandler(testSessSvc, mypageSvc, repo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://math-overcome.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("session", store))

	auth := r.Group("/auth")
	{
		auth.GET("/ping", authHandler.Ping)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/register", authHandler.Register)
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
