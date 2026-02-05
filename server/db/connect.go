package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"eefu/models"
)

func ConnectDB() *gorm.DB {

	// Ensure volume directory exists
	if err := os.MkdirAll("/data", 0755); err != nil {
		log.Fatal("failed to create data dir:", err)
	}

	var dbPath string
	if os.Getenv("FLY_ALLOC_ID") != "" {
		dbPath = "/data/eefu.db"
	} else {
		dbPath = "./data/eefu.db"
	}

	dsn := "file:" + dbPath + "?_journal_mode=WAL&_foreign_keys=on"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// For production, skip this and use Goose or Gormigrate migrations
	if err := db.AutoMigrate(
		&models.Genre{},
		&models.Tag{},
		&models.Singer{},
		&models.SingerAlias{},
		&models.Song{},
		&models.SongAlias{},
		&models.SongSinger{},
		&models.MediaSource{},
		&models.AudioMetadata{},
		&models.VideoMetadata{},
		&models.ImageMetadata{},
	); err != nil {
		log.Fatal("auto migration failed:", err)
	}

	return db
}
