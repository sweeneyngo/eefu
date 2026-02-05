package models

import (
	"time"

	"gorm.io/gorm"
)

type SongType string

const (
	SongTypeOriginal    SongType = "original"
	SongTypeCover       SongType = "cover"
	SongTypeRemix       SongType = "remix"
	SongTypeCompilation SongType = "compilation"
)

type Song struct {
	gorm.Model
	HashID          string   `gorm:"unique;not null" json:"hash_id"`
	Title           string   `gorm:"not null;size:100;uniqueIndex:idx_title_version" json:"title"`
	Description     string   `json:"description"`
	Type            SongType `gorm:"not null" json:"type"`
	OriginalSongID  *uint    `json:"-"`
	OriginalSong    *Song
	Version         int           `gorm:"default:1;uniqueIndex:idx_title_version" json:"version"`
	ReleasedAt      *time.Time    `json:"released_at"`
	Genres          []Genre       `gorm:"many2many:song_genres" json:"genres"`
	MediaSources    []MediaSource `gorm:"foreignKey:SongHashID;references:HashID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"media_sources"`
	Singers         []Singer      `gorm:"many2many:song_singers" json:"singers"`
	SongSingers     []SongSinger  `gorm:"foreignKey:SongID"`
	Tags            []Tag         `gorm:"many2many:song_tags" json:"tags"`
	Aliases         []SongAlias   `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE" json:"aliases"`
	SongGroupHashID string        `gorm:"index" json:"song_group_hash_id"`
}

type SongAlias struct {
	gorm.Model
	SongID   uint   `gorm:"not null;index:idx_song_alias_name,unique" json:"-"`
	Name     string `gorm:"not null;index:idx_song_alias_name,unique" json:"name"`
	Language string `gorm:"size:8"` // ISO code

	Song Song `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE" json:"-"`
}
