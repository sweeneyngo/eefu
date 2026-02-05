package serializers

import (
	"eefu/models"
	"time"
)

type SongAliasSerializer struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type SongSingerSerializer struct {
	Singer SingerSerializer  `json:"singer"`
	Role   models.SingerRole `json:"role"`
}

type SongSingerPublicSerializer struct {
	Singer SingerPublicSerializer `json:"singer"`
	Role   models.SingerRole      `json:"role"`
}

type SongPublicSerializer struct {
	HashID      string                        `json:"hash_id"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	Type        models.SongType               `json:"type"`
	ReleasedAt  *time.Time                    `json:"released_at"`
	Genres      []GenrePublicSerializer       `json:"genres"`
	Tags        []TagPublicSerializer         `json:"tags"`
	Singers     []SongSingerPublicSerializer  `json:"singers"`
	Media       []MediaSourcePublicSerializer `json:"media_sources"`
	Aliases     []SongAliasSerializer         `json:"aliases"`
	Version     int                           `json:"version"`
}

type SongSerializer struct {
	ID          uint                    `json:"id"`
	HashID      string                  `json:"hash_id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Type        models.SongType         `json:"type"`
	ReleasedAt  *time.Time              `json:"released_at"`
	Genres      []GenreSerializer       `json:"genres"`
	Tags        []TagSerializer         `json:"tags"`
	Singers     []SongSingerSerializer  `json:"singers"`
	Media       []MediaSourceSerializer `json:"media_sources"`
	Aliases     []SongAliasSerializer   `json:"aliases"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Version     int                     `json:"version"`
}

func SerializeSongPublic(song models.Song) SongPublicSerializer {
	genres := make([]GenrePublicSerializer, len(song.Genres))
	for i, g := range song.Genres {
		genres[i] = SerializeGenrePublic(g)
	}

	tags := make([]TagPublicSerializer, len(song.Tags))
	for i, t := range song.Tags {
		tags[i] = SerializeTagPublic(t)
	}

	singers := make([]SongSingerPublicSerializer, len(song.SongSingers))
	for i, ss := range song.SongSingers {
		singers[i] = SongSingerPublicSerializer{
			Singer: SerializeSingerPublic(ss.Singer),
			Role:   ss.Role,
		}
	}
	media := make([]MediaSourcePublicSerializer, len(song.MediaSources))
	for i, m := range song.MediaSources {
		media[i] = SerializeMediaSourcePublic(m)
	}

	aliases := make([]SongAliasSerializer, len(song.Aliases))
	for i, a := range song.Aliases {
		aliases[i] = SongAliasSerializer{
			Name:     a.Name,
			Language: a.Language,
		}
	}

	return SongPublicSerializer{
		HashID:      song.HashID,
		Title:       song.Title,
		Description: song.Description,
		Type:        song.Type,
		Version:     song.Version,
		Genres:      genres,
		Tags:        tags,
		Singers:     singers,
		Media:       media,
		Aliases:     aliases,
		ReleasedAt:  song.ReleasedAt,
	}
}

func SerializeSong(song models.Song) SongSerializer {
	genres := make([]GenreSerializer, len(song.Genres))
	for i, g := range song.Genres {
		genres[i] = SerializeGenre(g)
	}

	tags := make([]TagSerializer, len(song.Tags))
	for i, t := range song.Tags {
		tags[i] = SerializeTag(t)
	}

	singers := make([]SongSingerSerializer, len(song.SongSingers))
	for i, ss := range song.SongSingers {
		singers[i] = SongSingerSerializer{
			Singer: SerializeSinger(ss.Singer),
			Role:   ss.Role,
		}
	}

	media := make([]MediaSourceSerializer, len(song.MediaSources))
	for i, m := range song.MediaSources {
		media[i] = SerializeMediaSource(m)
	}

	aliases := make([]SongAliasSerializer, len(song.Aliases))
	for i, a := range song.Aliases {
		aliases[i] = SongAliasSerializer{
			Name:     a.Name,
			Language: a.Language,
		}
	}

	return SongSerializer{
		ID:          song.ID,
		HashID:      song.HashID,
		Title:       song.Title,
		Description: song.Description,
		Type:        song.Type,
		Genres:      genres,
		Tags:        tags,
		Singers:     singers,
		Media:       media,
		Aliases:     aliases,
		CreatedAt:   song.CreatedAt,
		UpdatedAt:   song.UpdatedAt,
		Version:     song.Version,
		ReleasedAt:  song.ReleasedAt,
	}
}
