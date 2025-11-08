package main

import (
	"log"

	"app/internal/auth"
	"app/internal/db"
	"app/internal/user"

	"github.com/gin-contrib/cors"
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

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(googleOAuth, jwtManager, userRepo)
	authHandler := auth.NewHandler(authService)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		// Auth
		protected.GET("/auth/me", authHandler.Me)
		protected.GET("/auth/logout", authHandler.Logout)

		// User
		protected.DELETE("/user", userHandler.DeleteCurrentUser)
		protected.PUT("/user", userHandler.UpdateCurrentUser)
	}

	router.Run(":8080")
}
