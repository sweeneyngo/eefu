package main

import (
	"log"
	"net/http"
	"os"

	"eefu/db"
	"eefu/routes"
	"eefu/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env")
	}

	endpoint := os.Getenv("CLOUDFLARE_DEFAULT_ENDPOINT")
	accessKey := os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")
	accessSecret := os.Getenv("CLOUDFLARE_SECRET_ACCESS_KEY")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	database := db.ConnectDB()

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("failed to get underlying sql.DB:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("database ping failed:", err)
	}
	log.Println("Database connected successfully")

	s3 := storage.NewR2Client(endpoint, accessKey, accessSecret)
	presigner := storage.NewPresigner(s3)
	uploader := storage.NewR2Uploader(s3)

	routes.RegisterRoutes(r, database, s3, presigner, uploader)
	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
