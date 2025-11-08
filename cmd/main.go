package main

import (
	"log"

	"app/internal/auth"
	"app/internal/cv"
	"app/internal/db"
	"app/internal/tag"
	"app/internal/test"
	"app/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("failed connect to db: ", err)
	}
	defer db.Close()

	
	googleOAuth, err := auth.NewGoogleOAuth()
	if err != nil {
		log.Fatal(err)
	}
	
	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatal(err)
	}
	
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	authService := auth.NewAuthService(googleOAuth, jwtManager, userRepo)
	authHandler := auth.NewAuthHandler(authService)
	
	cvRepo := cv.NewCVRepository(db)
	cvService := cv.NewCVService(cvRepo)
	cvHandler := cv.NewCVHandler(cvService)

	tagRepo := tag.NewTagRepository(db)
	tagService := tag.NewTagService(tagRepo)
	tagHandler := tag.NewTagHandler(tagService)

	testRepo := test.NewTestRepository(db)
	testService := test.NewTestService(testRepo)
	testHandler := test.NewTestHandler(testService)

	router := gin.Default()
	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	router.GET("/tags", tagHandler.GetAllTags)
	router.GET("/tags/:tag_id", tagHandler.GetTagById)

	router.GET("/test", testHandler.GetTestsByTags)
	router.GET("test/:test_id", testHandler.GetTestById)
	router.POST("/test", testHandler.CreateTest)
	router.DELETE("/test/:test_id", testHandler.DeleteTest)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", authHandler.Me)
		protected.DELETE("/user", userHandler.DeleteCurrentUser)
		protected.PUT("/user", userHandler.UpdateCurrentUser)

		protected.GET("/cvs", cvHandler.GetAllCVsByCurrentUser)
		protected.GET("/cvs/:cv_id", cvHandler.GetCVByID)
		protected.POST("/cvs", cvHandler.CreateCV)
		protected.PUT("/cvs/:cv_id", cvHandler.UpdateCV)
		protected.DELETE("/cvs/:cv_id", cvHandler.DeleteCV)
	}

	router.Run(":8080")
}
