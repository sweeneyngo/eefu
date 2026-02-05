package services

import (
	"context"
	"eefu/models"
	"eefu/storage"
	"eefu/types"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const CDN_URL = "https://cdn.ifuxyl.dev"

type SongService struct {
	DB        *gorm.DB
	Presigner *storage.Presigner
}

func NewSongService(db *gorm.DB, presigner *storage.Presigner) *SongService {
	return &SongService{DB: db, Presigner: presigner}
}

func (s *SongService) GetAll(ctx context.Context) ([]models.Song, error) {
	var songs []models.Song
	err := s.DB.WithContext(ctx).
		Preload("Genres").
		Preload("MediaSources").
		Preload("SongSingers.Singer").
		Preload("Tags").
		Preload("Aliases").
		Find(&songs).Error
	return songs, err
}

func (s *SongService) GetVersionsByGroup(ctx context.Context, songGroupHashID string) ([]models.Song, error) {
	var versions []models.Song
	err := s.DB.WithContext(ctx).
		Preload("Genres").
		Preload("Tags").
		Preload("Aliases").
		Preload("SongSingers.Singer").
		Preload("MediaSources").
		Preload("MediaSources.AudioMetadata").
		Preload("MediaSources.VideoMetadata").
		Preload("MediaSources.ImageMetadata").
		Where("song_group_hash_id = ?", songGroupHashID).
		Order("version DESC").
		Find(&versions).Error
	return versions, err
}

func (s *SongService) GetMedia(ctx context.Context, hashID string) ([]types.MediaMinimal, error) {
	var song models.Song
	if err := s.DB.WithContext(ctx).Preload("MediaSources").First(&song, "hash_id = ?", hashID).Error; err != nil {
		return nil, err
	}

	mediaList := make([]types.MediaMinimal, 0, len(song.MediaSources))
	for _, media := range song.MediaSources {
		var url string
		isArt := media.FileType == models.MediaSourceFileTypeArt
		isMP3 := media.FormatType == models.MediaSourceFormatTypeMP3

		if isArt || isMP3 {
			url = fmt.Sprintf("%s/%s", CDN_URL, media.URL)
		} else {
			var err error
			url, err = s.Presigner.GeneratePresignedURL(ctx, BUCKET_NAME, media.URL)
			if err != nil {
				continue
			}
		}

		mediaList = append(mediaList, types.MediaMinimal{
			FileType: media.FileType,
			URL:      url,
		})
	}

	return mediaList, nil
}

func (s *SongService) CreateSong(ctx context.Context, input types.SongInput) (models.Song, error) {
	var song models.Song

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error

		song = buildSong(input)

		if err = s.checkDuplicate(tx, input.Title); err != nil {
			return err
		}

		song.Version, err = s.nextVersion(tx, input.Title)
		if err != nil {
			return err
		}

		if err = tx.Create(&song).Error; err != nil {
			return err
		}

		if err = s.assignExistingGenres(tx, &song, input.Genres); err != nil {
			return err
		}
		if err = s.assignExistingSingers(tx, &song, input.Singers); err != nil {
			return err
		}
		if err = s.assignExistingTags(tx, &song, input.Tags); err != nil {
			return err
		}
		if err = s.addAliases(tx, &song, input.Aliases); err != nil {
			return err
		}

		return tx.Preload("Genres").
			Preload("SongSingers.Singer").
			Preload("Tags").
			Preload("MediaSources").
			Preload("Aliases").
			First(&song, song.ID).Error
	})

	return song, err
}

func (s *SongService) CreateSongVersion(ctx context.Context, hashID string, input types.SongVersionInput) (*models.Song, error) {
	var newSong models.Song
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var currentSong models.Song
		var err error

		currentSong, err = s.getSong(tx, hashID)
		if err != nil {
			return err
		}

		var maxVersion int
		if err := tx.Model(&models.Song{}).
			Where("song_group_hash_id = ?", currentSong.SongGroupHashID).
			Select("COALESCE(MAX(version), 0)").
			Scan(&maxVersion).Error; err != nil {
			return err
		}

		newSong = cloneSong(currentSong, input)
		newSong.Version = maxVersion + 1
		if err := tx.Create(&newSong).Error; err != nil {
			return err
		}

		if len(input.Tags) > 0 {
			if err := s.assignExistingTags(tx, &newSong, input.Tags); err != nil {
				return err
			}
		} else {
			if err := tx.Model(&newSong).Association("Tags").Replace(currentSong.Tags); err != nil {
				return err
			}
		}

		if len(input.Singers) > 0 {
			if err := s.assignExistingSingers(tx, &newSong, input.Singers); err != nil {
				return err
			}
		} else {
			if err := s.copySingers(tx, &newSong, currentSong.SongSingers); err != nil {
				return err
			}
		}

		if err := tx.Model(&newSong).Association("Genres").Replace(currentSong.Genres); err != nil {
			return err
		}

		if err := tx.Model(&newSong).Association("Aliases").Replace(currentSong.Aliases); err != nil {
			return err
		}

		return tx.Preload("Genres").
			Preload("SongSingers.Singer").
			Preload("Tags").
			Preload("MediaSources").
			Preload("Aliases").
			First(&newSong, newSong.ID).Error
	})
	return &newSong, err
}

func buildSong(input types.SongInput) models.Song {
	song := models.Song{
		Title:           input.Title,
		Type:            input.Type,
		HashID:          uuid.New().String(),
		SongGroupHashID: uuid.New().String(),
	}
	if input.Description != nil {
		song.Description = *input.Description
	}
	if input.ReleasedAt != nil {
		song.ReleasedAt = input.ReleasedAt
	}
	return song
}

func cloneSong(prevSong models.Song, input types.SongVersionInput) models.Song {
	song := models.Song{
		Title:           prevSong.Title,
		Type:            prevSong.Type,
		HashID:          uuid.New().String(),
		SongGroupHashID: prevSong.SongGroupHashID,
	}
	if input.Description != nil {
		song.Description = *input.Description
	}
	if input.ReleasedAt != nil {
		song.ReleasedAt = input.ReleasedAt
	}
	return song
}

func (s *SongService) getSong(tx *gorm.DB, hashID string) (models.Song, error) {
	var song models.Song
	err := tx.Model(&models.Song{}).
		Preload("Genres").
		Preload("MediaSources").
		Preload("SongSingers.Singer").
		Preload("Tags").
		Where("hash_id = ?", hashID).
		First(&song).Error
	return song, err
}

func (s *SongService) checkDuplicate(tx *gorm.DB, title string) error {
	var existing int64
	if err := tx.Model(&models.Song{}).Where("title = ?", title).Count(&existing).Error; err != nil {
		return err
	}
	if existing > 0 {
		return errors.New("song with this title already exists")
	}
	return nil
}

func (s *SongService) nextVersion(tx *gorm.DB, title string) (int, error) {
	var maxVersion int
	if err := tx.Model(&models.Song{}).
		Where("title = ?", title).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error; err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

func (s *SongService) assignExistingGenres(tx *gorm.DB, song *models.Song, genres []types.GenreInput) error {
	if len(genres) == 0 {
		return nil
	}
	names := make([]string, len(genres))
	for i, g := range genres {
		names[i] = g.Name
	}
	var found []models.Genre
	if err := tx.Where("name IN ?", names).Find(&found).Error; err != nil {
		return err
	}
	return tx.Model(song).Association("Genres").Replace(found)
}

func (s *SongService) assignExistingSingers(tx *gorm.DB, song *models.Song, singers []types.SongSingerInput) error {
	if len(singers) == 0 {
		return nil
	}
	var songSingers []models.SongSinger
	for _, s := range singers {
		var singer models.Singer
		if err := tx.Where("name = ?", s.Name).First(&singer).Error; err != nil {
			return err
		}
		songSingers = append(songSingers, models.SongSinger{
			SongID:   song.ID,
			SingerID: singer.ID,
			Role:     models.SingerRole(s.Role),
		})
	}
	return tx.Create(&songSingers).Error
}

func (s *SongService) copySingers(tx *gorm.DB, newSong *models.Song, prevSingers []models.SongSinger) error {
	if len(prevSingers) == 0 {
		return nil
	}

	songSingers := make([]models.SongSinger, len(prevSingers))
	for i, ps := range prevSingers {
		songSingers[i] = models.SongSinger{
			SongID:   newSong.ID,
			SingerID: ps.SingerID,
			Role:     ps.Role,
		}
	}
	return tx.Create(&songSingers).Error
}

func (s *SongService) assignExistingTags(tx *gorm.DB, song *models.Song, tags []types.TagInput) error {
	if len(tags) == 0 {
		return nil
	}
	var found []models.Tag
	for _, t := range tags {
		var tag models.Tag
		if err := tx.Where("name = ? AND type = ?", t.Name, t.Type).First(&tag).Error; err != nil {
			return err
		}
		found = append(found, tag)
	}
	return tx.Model(song).Association("Tags").Replace(found)
}

func (s *SongService) addAliases(tx *gorm.DB, song *models.Song, aliases []types.SongAliasInput) error {
	if len(aliases) == 0 {
		return nil
	}
	var toCreate []models.SongAlias
	for _, a := range aliases {
		toCreate = append(toCreate, models.SongAlias{
			SongID:   song.ID,
			Name:     a.Name,
			Language: a.Language,
		})
	}
	return tx.Create(&toCreate).Error
}
