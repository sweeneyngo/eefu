package serializers

import (
	"eefu/models"
)

type GenrePublicSerializer struct {
	Name string `json:"name"`
}

type GenreSerializer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func SerializeGenrePublic(genre models.Genre) GenrePublicSerializer {
	return GenrePublicSerializer{
		Name: genre.Name,
	}
}

func SerializeGenre(genre models.Genre) GenreSerializer {
	return GenreSerializer{
		ID:   genre.ID,
		Name: genre.Name,
	}
}
