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

	userRepo := user.NewUserRepository(db)

	googleOAuth, err := auth.NewGoogleOAuth()
	if err != nil {
		log.Fatal(err)
	}
	
	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatal(err)
	}
	
	authService := auth.NewAuthService(googleOAuth, jwtManager, userRepo)
	authHandler := auth.NewAuthHandler(authService)
	

	router := gin.Default()
	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", authHandler.Me)
	}

	router.Run(":8080")
}
