package models

import (
	"gorm.io/gorm"
)

const GenreNameMaxLength = 32

type Genre struct {
	gorm.Model
	Name string `gorm:"unique;not null;size:32" json:"name"`
}
