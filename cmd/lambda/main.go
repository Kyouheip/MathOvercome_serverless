package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/handler"
	_ "github.com/Kyouheip/MathOvercome_serverless/internal/model" // gob.Register の init() を呼ぶ
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// DB接続
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// DI
	repo := repository.NewRepository(db)
	loginSvc := service.NewLoginService(repo)
	testSessSvc := service.NewTestSessionService(repo)
	mypageSvc := service.NewMypageService(repo)
	authHandler := handler.NewAuthHandler(loginSvc)
	sessionHandler := handler.NewSessionHandler(testSessSvc, mypageSvc, repo)

	r := gin.Default()

	// CORS (Java の @CrossOrigin と同じ設定)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:3000",
			"https://math-overcome.vercel.app",
			"http://52.68.88.3",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// セッション (Cookie Store)
	// TODO: Lambda 運用時は DynamoDB や Redis ベースのストアに切り替える
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("session", store))

	// ルーティング (Java の @RequestMapping と一致)
	auth := r.Group("/auth")
	{
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

	// TODO: Lambda デプロイ時は以下に切り替え
	// ginLambda := ginadapter.New(r)
	// lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 	return ginLambda.ProxyWithContext(ctx, req)
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
