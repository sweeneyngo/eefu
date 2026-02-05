package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"eefu/models"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

const BUCKET_NAME = "eefu-media"
const BUCKET_CDN_NAME = "eefu-cdn"

type MediaService struct {
	DB       *gorm.DB
	S3Client *s3.Client
	Uploader *manager.Uploader
	Bucket   string
}

func NewMediaService(db *gorm.DB, s3Client *s3.Client, uploader *manager.Uploader) *MediaService {
	return &MediaService{DB: db, S3Client: s3Client, Uploader: uploader, Bucket: BUCKET_NAME}
}
func (s *MediaService) UploadMedia(ctx context.Context, songHashID string, file io.Reader, fileName string, fileType models.MediaSourceFileType, storageType models.MediaSourceStorageType) (*models.MediaSource, error) {
	song, err := s.getSong(songHashID)
	if err != nil {
		return nil, err
	}

	data, err := readFile(file)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	formatType, ok := models.ExtToFormat[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	bucket, objectKey := getTarget(fileName, song.HashID, fileType, formatType)
	if err := s.uploadFile(ctx, bucket, objectKey, data); err != nil {
		return nil, err
	}

	var media *models.MediaSource
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		media, err = s.saveMediaWithMetadata(tx, song.HashID, data, objectKey, fileType, formatType, storageType)
		return err
	})

	if err != nil {
		s.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		})
		return nil, err
	}

	return media, nil
}

func (s *MediaService) getSong(hashID string) (*models.Song, error) {
	var song models.Song
	if err := s.DB.First(&song, "hash_id = ?", hashID).Error; err != nil {
		return nil, fmt.Errorf("song not found: %w", err)
	}
	return &song, nil
}

func readFile(file io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	data := buf.Bytes()
	return data, nil
}

func getTarget(fileName, songHashID string, fileType models.MediaSourceFileType, formatType models.MediaSourceFormatType) (bucket string, objectKey string) {
	bucket = BUCKET_NAME
	if fileType == models.MediaSourceFileTypeArt || formatType == models.MediaSourceFormatTypeMP3 {
		bucket = BUCKET_CDN_NAME
	}
	objectKey = fmt.Sprintf("%s/%d-%s", songHashID, time.Now().UnixNano(), fileName)
	return bucket, objectKey
}

func (s *MediaService) uploadFile(ctx context.Context, bucket, objectKey string, data []byte) error {
	_, err := s.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *MediaService) saveMediaWithMetadata(tx *gorm.DB, songHashID string, data []byte, objectKey string, fileType models.MediaSourceFileType, formatType models.MediaSourceFormatType, storageType models.MediaSourceStorageType) (*models.MediaSource, error) {
	media := models.MediaSource{
		SongHashID:  songHashID,
		URL:         objectKey,
		StorageType: storageType,
		FileType:    fileType,
		FormatType:  formatType,
		Checksum:    fmt.Sprintf("%x", sha256.Sum256(data)),
	}

	if err := tx.Create(&media).Error; err != nil {
		return nil, err
	}

	var meta interface{}
	var err error

	switch fileType {
	case models.MediaSourceFileTypeAudio:
		meta, err = extractAudioMetadata(data)
	case models.MediaSourceFileTypeVideo:
		meta, err = extractVideoMetadata(data)
	case models.MediaSourceFileTypeArt:
		meta, err = extractImageMetadata(data)
	}

	if err != nil {
		return nil, err
	}

	if meta != nil {
		switch m := meta.(type) {
		case *models.AudioMetadata:
			m.MediaSourceID = media.ID
			if err := tx.Create(m).Error; err != nil {
				return nil, err
			}
		case *models.VideoMetadata:
			m.MediaSourceID = media.ID
			if err := tx.Create(m).Error; err != nil {
				return nil, err
			}
		case *models.ImageMetadata:
			m.MediaSourceID = media.ID
			if err := tx.Create(m).Error; err != nil {
				return nil, err
			}
		}
	}

	return &media, nil
}

func extractAudioMetadata(data []byte) (*models.AudioMetadata, error) {
	tmpFile, err := os.CreateTemp("", "audio-*")
	if err != nil {
		return &models.AudioMetadata{}, err
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(data); err != nil {
		return &models.AudioMetadata{}, err
	}
	tmpFile.Close()

	// Use ffprobe to get audio info as JSON
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", tmpFile.Name())
	out, err := cmd.Output()
	if err != nil {
		return &models.AudioMetadata{}, err
	}

	var ffprobe struct {
		Streams []struct {
			CodecType     string `json:"codec_type"`
			Channels      int    `json:"channels"`
			SampleRate    string `json:"sample_rate"`
			BitsPerSample int    `json:"bits_per_sample"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	if err := json.Unmarshal(out, &ffprobe); err != nil {
		return &models.AudioMetadata{}, err
	}

	var stream struct {
		Channels      int
		SampleRate    int
		BitsPerSample int
	}
	for _, s := range ffprobe.Streams {
		if s.CodecType == "audio" {
			stream.Channels = s.Channels
			stream.SampleRate, _ = strconv.Atoi(s.SampleRate)
			stream.BitsPerSample = s.BitsPerSample
			break
		}
	}

	duration, _ := strconv.ParseFloat(ffprobe.Format.Duration, 64)

	return &models.AudioMetadata{
		Bitrate:       0,
		Channels:      stream.Channels,
		SampleRate:    stream.SampleRate,
		BitsPerSample: stream.BitsPerSample,
		Duration:      duration,
	}, nil
}

func extractVideoMetadata(data []byte) (*models.VideoMetadata, error) {
	tmpFile, err := os.CreateTemp("", "video-*")
	if err != nil {
		return &models.VideoMetadata{}, err
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(data); err != nil {
		return &models.VideoMetadata{}, err
	}
	tmpFile.Close()

	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", "-show_format", tmpFile.Name())
	out, err := cmd.Output()
	if err != nil {
		return &models.VideoMetadata{}, err
	}

	var ffprobe struct {
		Streams []struct {
			CodecType  string `json:"codec_type"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			RFrameRate string `json:"r_frame_rate"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	if err := json.Unmarshal(out, &ffprobe); err != nil {
		return &models.VideoMetadata{}, err
	}

	var stream struct {
		Width     int
		Height    int
		FrameRate int
	}
	for _, s := range ffprobe.Streams {
		if s.CodecType == "video" {
			stream.Width = s.Width
			stream.Height = s.Height
			parts := strings.Split(s.RFrameRate, "/")
			if len(parts) == 2 {
				num, _ := strconv.Atoi(parts[0])
				den, _ := strconv.Atoi(parts[1])
				if den != 0 {
					stream.FrameRate = num / den
				}
			}
			break
		}
	}

	duration, _ := strconv.ParseFloat(ffprobe.Format.Duration, 64)

	return &models.VideoMetadata{
		Width:     stream.Width,
		Height:    stream.Height,
		FrameRate: stream.FrameRate,
		Duration:  duration,
	}, nil
}

func extractImageMetadata(data []byte) (*models.ImageMetadata, error) {
	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return &models.ImageMetadata{}, err
	}

	return &models.ImageMetadata{
		Width:  img.Width,
		Height: img.Height,
	}, nil
}
