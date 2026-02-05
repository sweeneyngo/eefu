package routes

import (
	"eefu/handlers"
	"eefu/middleware"
	"eefu/services"
	"eefu/storage"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func RegisterRoutes(r chi.Router, db *gorm.DB, s3 *s3.Client, presigner *storage.Presigner, uploader *manager.Uploader) {

	songService := services.NewSongService(db, presigner)
	mediaService := services.NewMediaService(db, s3, uploader)

	// Public endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.OptionalAdminMiddleware(os.Getenv("API_SECRET_KEY")))
		r.Use(middleware.WithRequestLogger)
		r.Use(middleware.RequestID)

		// Songs
		r.Route("/songs", func(r chi.Router) {
			r.Get("/", handlers.GetSongs(songService))
			r.Get("/group/{song_group_hash_id}/versions", handlers.GetSongVersionsByGroup(songService))
			r.Get("/{hash_id}/download", handlers.GetSongMedia(songService))
		})
	})

	// Admin-only endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.APIKeyAuthMiddleware(os.Getenv("API_SECRET_KEY")))
		r.Use(middleware.WithRequestLogger)
		r.Use(middleware.RequestID)

		// Songs
		r.Route("/admin/songs", func(r chi.Router) {
			r.Get("/", handlers.GetSongs(songService))
			r.Post("/", handlers.CreateSong(songService))
			r.Post("/{hash_id}/version", handlers.CreateSongVersion(songService))
			r.Post("/{hash_id}/upload", handlers.UploadMedia(mediaService))
		})

		// Genres
		r.Route("/admin/genres", func(r chi.Router) {
			r.Get("/", handlers.GetGenres(db))
			r.Post("/", handlers.CreateGenre(db))
		})

		// Tags
		r.Route("/admin/tags", func(r chi.Router) {
			r.Get("/", handlers.GetTags(db))
			r.Post("/", handlers.CreateTag(db))
		})

		// Singers
		r.Route("/admin/singers", func(r chi.Router) {
			r.Get("/", handlers.GetSingers(db))
			r.Post("/", handlers.CreateSinger(db))
		})
	})
}
