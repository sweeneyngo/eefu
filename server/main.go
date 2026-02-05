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
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, rely on Fly variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	endpoint := os.Getenv("CLOUDFLARE_DEFAULT_ENDPOINT")
	accessKey := os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")
	accessSecret := os.Getenv("CLOUDFLARE_SECRET_ACCESS_KEY")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://eefu-gamma.vercel.app/", "https://www.ifuxyl.dev"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

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

	addr := "0.0.0.0:" + port
	log.Println("Server starting on", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("server failed:", err)
	}
}
