package models

import (
	"gorm.io/gorm"
)

type TagType string

const (
	TagTypeArchive TagType = "archive"
	TagTypeCustom  TagType = "custom"
	TagTypeStage   TagType = "stage"
	TagTypeStatus  TagType = "status"

	TagNameMaxLength        = 32
	TagDescriptionMaxLength = 256
)

type Tag struct {
	gorm.Model
	Name        string  `gorm:"not null;size:32;uniqueIndex:uniq_name_type" json:"name"`
	Type        TagType `gorm:"not null;default:custom;uniqueIndex:uniq_name_type" json:"type"`
	Description string  `gorm:"size:256" json:"description"`
}
