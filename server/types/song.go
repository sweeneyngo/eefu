package types

import (
	"eefu/models"
	"time"
)

type SongInput struct {
	Title       string          `json:"title" validate:"required,min=1,max=128"`
	Description *string         `json:"description" validate:"omitempty,max=500"`
	Type        models.SongType `json:"type" validate:"required,oneof=original cover remix compilation"`
	ReleasedAt  *time.Time      `json:"released_at,omitempty"`

	Genres  []GenreInput      `json:"genres" validate:"omitempty,dive"`
	Singers []SongSingerInput `json:"singers" validate:"omitempty,dive"`
	Tags    []TagInput        `json:"tags" validate:"omitempty,dive"`
	Aliases []SongAliasInput  `json:"aliases" validate:"omitempty,dive"`
}

// Specifies what can be overwritten.
// When overriding, only overwrite with these fields if they're included.
type SongVersionInput struct {
	Tags        []TagInput        `json:"tags,omitempty"`
	Singers     []SongSingerInput `json:"singers,omitempty"`
	Description *string           `json:"description,omitempty"`
	ReleasedAt  *time.Time        `json:"released_at,omitempty"`
}

type SongAliasInput struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Language string `json:"language" validate:"required,len=2"` // ISO 639-1 code
}

type SongSingerInput struct {
	Name string `json:"name" validate:"required,min=1,max=64"`
	Role string `json:"role" validate:"required,oneof=main featured"`
}
