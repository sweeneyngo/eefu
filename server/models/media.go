package models

import "gorm.io/gorm"

type MediaSourceStorageType string
type MediaSourceFileType string
type MediaSourceFormatType string

const (
	MediaSourceStorageTypeLocal      MediaSourceStorageType = "local"
	MediaSourceStorageTypeYouTube    MediaSourceStorageType = "youtube"
	MediaSourceStorageTypeS3         MediaSourceStorageType = "s3"
	MediaSourceStorageTypeSoundCloud MediaSourceStorageType = "soundcloud"
)

const (
	MediaSourceFileTypeArt   MediaSourceFileType = "art"
	MediaSourceFileTypeAudio MediaSourceFileType = "audio"
	MediaSourceFileTypeVideo MediaSourceFileType = "video"
)

const (
	MediaSourceFormatTypeMP3  MediaSourceFormatType = "mp3"
	MediaSourceFormatTypeWAV  MediaSourceFormatType = "wav"
	MediaSourceFormatTypeFLAC MediaSourceFormatType = "flac"
	MediaSourceFormatTypeMP4  MediaSourceFormatType = "mp4"
	MediaSourceFormatTypeWEBM MediaSourceFormatType = "webm"
	MediaSourceFormatTypePNG  MediaSourceFormatType = "png"
	MediaSourceFormatTypeJPG  MediaSourceFormatType = "jpg"
	MediaSourceFormatTypeJPEG MediaSourceFormatType = "jpeg"
	MediaSourceFormatTypeWEBP MediaSourceFormatType = "webp"
)

var ExtToFormat = map[string]MediaSourceFormatType{
	".mp3":  MediaSourceFormatTypeMP3,
	".wav":  MediaSourceFormatTypeWAV,
	".flac": MediaSourceFormatTypeFLAC,
	".mp4":  MediaSourceFormatTypeMP4,
	".webm": MediaSourceFormatTypeWEBM,
	".png":  MediaSourceFormatTypePNG,
	".jpg":  MediaSourceFormatTypeJPG,
	".jpeg": MediaSourceFormatTypeJPEG,
	".webp": MediaSourceFormatTypeWEBP,
}

type MediaSource struct {
	gorm.Model
	SongHashID  string                 `gorm:"not null;index" json:"-"`
	StorageType MediaSourceStorageType `gorm:"not null" json:"storage_type"`
	URL         string                 `gorm:"not null;uniqueIndex" json:"url"`
	FileType    MediaSourceFileType    `gorm:"default:'audio'" json:"file_type"`
	FormatType  MediaSourceFormatType  `gorm:"not null" json:"format_type"`
	Checksum    string                 `gorm:"size:64;index" json:"checksum"`

	Song Song `gorm:"foreignKey:SongHashID;references:HashID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`

	AudioMetadata *AudioMetadata `gorm:"foreignKey:MediaSourceID" json:"audio_metadata,omitempty"`
	VideoMetadata *VideoMetadata `gorm:"foreignKey:MediaSourceID" json:"video_metadata,omitempty"`
	ImageMetadata *ImageMetadata `gorm:"foreignKey:MediaSourceID" json:"image_metadata,omitempty"`
}

type AudioMetadata struct {
	gorm.Model
	MediaSourceID uint    `gorm:"uniqueIndex" json:"-"`
	Bitrate       int     `gorm:"not null" json:"bitrate"`
	SampleRate    int     `gorm:"not null" json:"sample_rate"`
	Channels      int     `gorm:"not null" json:"channels"`
	BitsPerSample int     `gorm:"not null" json:"bits_per_sample"`
	Duration      float64 `gorm:"not null" json:"duration"`
}

type VideoMetadata struct {
	gorm.Model
	MediaSourceID uint    `gorm:"uniqueIndex" json:"-"`
	Width         int     `gorm:"not null" json:"width"`
	Height        int     `gorm:"not null" json:"height"`
	FrameRate     int     `gorm:"not null" json:"frame_rate"`
	Duration      float64 `gorm:"not null" json:"duration"`
}

type ImageMetadata struct {
	MediaSourceID uint `gorm:"uniqueIndex" json:"-"`
	Width         int  `gorm:"not null" json:"width"`
	Height        int  `gorm:"not null" json:"height"`
}
