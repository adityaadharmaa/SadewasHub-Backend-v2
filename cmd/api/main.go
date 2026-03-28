package main

import (
	"log"
	"os"
	"sadewashub-go/internal/config"
	"sadewashub-go/internal/controllers"
	"sadewashub-go/internal/routes"
	"sadewashub-go/internal/seeders"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠ Warning File .env tidak ditemukan")
	}

	config.ConnectionDatabase()

	seeders.RunSeeder()

	controllers.InitOAuth()

	r := gin.Default()

	v2 := r.Group("/api/v2")
	routes.SetupAuthRoutes(v2)

	// r.GET("/api/v2/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"status":  "success",
	// 		"message": "Pong! API SadewasHub (GOLANG) sudah berjalan 🚀",
	// 	})
	// })

	// authGroup := r.Group("/api/v2/auth")
	// {
	// 	authGroup.GET("/google", controllers.GoogleLogin)
	// 	authGroup.GET("/google/callback", controllers.GoogleCallback)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🔥 Server berjalan di http://localhost:%s\n", port)
	r.Run(":" + port)
}
