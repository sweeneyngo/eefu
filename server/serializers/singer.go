package serializers

import (
	"eefu/models"
)

type AliasDetail struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type SingerPublicSerializer struct {
	Name    string        `json:"name"`
	Aliases []AliasDetail `json:"aliases"`
}

type SingerSerializer struct {
	ID      uint          `json:"id"`
	Name    string        `json:"name"`
	Aliases []AliasDetail `json:"aliases"`
}

func SerializeSingerPublic(singer models.Singer) SingerPublicSerializer {
	aliases := make([]AliasDetail, len(singer.Aliases))
	for i, a := range singer.Aliases {
		aliases[i] = AliasDetail{
			Name:     a.Name,
			Language: a.Language,
		}
	}
	return SingerPublicSerializer{
		Name:    singer.Name,
		Aliases: aliases,
	}
}

func SerializeSinger(singer models.Singer) SingerSerializer {
	aliases := make([]AliasDetail, len(singer.Aliases))
	for i, a := range singer.Aliases {
		aliases[i] = AliasDetail{
			Name:     a.Name,
			Language: a.Language,
		}
	}
	return SingerSerializer{
		ID:      singer.ID,
		Name:    singer.Name,
		Aliases: aliases,
	}
}
