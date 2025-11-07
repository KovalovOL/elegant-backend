package main

import (
	"log"

	"app/internal/auth"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"app/internal/db"
	"app/internal/user"
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
	

	router := gin.Default()
	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

//	router.DELETE("/user", userHandler.DeleteCurrentUser)
//	router.PUT("/user", userHandler.UpdateCurrentUser)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", authHandler.Me)
		protected.DELETE("/user", userHandler.DeleteCurrentUser)
		protected.PUT("/user", userHandler.UpdateCurrentUser)
	}

	router.Run(":8080")
}
