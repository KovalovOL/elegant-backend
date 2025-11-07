package main

import (
	"log"

	"app/internal/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	googleOAuth, err := auth.NewGoogleOAuth()
	if err != nil {
		log.Fatal(err)
	}

	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewAuthService(googleOAuth, jwtManager)
	authHandler := auth.NewAuthHandler(authService)

	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", authHandler.Me)
	}

	router.Run(":8080")
}
