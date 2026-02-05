package models

import (
	"gorm.io/gorm"
)

type SingerRole string

const (
	SingerRoleMain            SingerRole = "main"
	SingerRoleFeatured        SingerRole = "featured"
	SingerNameMaxLength                  = 64
	SingerAliasNameMaxLength             = 64
	SingerAliasLanguageLength            = 2
)

type Singer struct {
	gorm.Model
	Name    string        `gorm:"unique;not null;size:64" json:"name"`
	Aliases []SingerAlias `gorm:"foreignKey:SingerID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"aliases"`
}

type SingerAlias struct {
	gorm.Model
	SingerID uint   `gorm:"not null;index" json:"-"`
	Name     string `gorm:"unique;not null;size:64" json:"name"`
	Language string `gorm:"not null;size:2" json:"language"` // ISO 639-1 codes

	Singer Singer `gorm:"foreignKey:SingerID;constraint:OnDelete:CASCADE,OnDelete:CASCADE" json:"-"`
}

type SongSinger struct {
	gorm.Model
	SongID   uint       `gorm:"not null;index:uniq_song_singer,unique" json:"song_id"`
	SingerID uint       `gorm:"not null;index:uniq_song_singer,unique" json:"singer_id"`
	Role     SingerRole `gorm:"not null" json:"role"`

	Song   Song   `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE,OnDelete:CASCADE" json:"-"`
	Singer Singer `gorm:"foreignKey:SingerID;constraint:OnDelete:CASCADE,OnDelete:CASCADE" json:"-"`
}
