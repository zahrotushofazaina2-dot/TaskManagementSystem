package main

import (
	"log"
	"os"

	"taskmanagement/config"
	"taskmanagement/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Mengatur mode Gin.
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Menghubungkan aplikasi ke database.
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Memuat konfigurasi autentikasi.
	authCfg, err := config.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Auth configuration failed: %v", err)
	}

	// Membuat router.
	router := routes.SetupRouter(db, authCfg)

	// Mengambil port dari environment.
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Menjalankan server.
	log.Printf("Server running on http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
