package main

import (
	"github.com/gin-gonic/gin"
	"log"

	"app/internal/auth"
	"app/internal/service"
	"app/internal/handler"
	"app/internal/middleware"
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

	authService := service.NewAuthService(googleOAuth, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", authHandler.Me)
	}

	router.Run(":8080")
}
