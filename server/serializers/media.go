package serializers

import (
	"eefu/models"
)

type AudioMetadataSerializer struct {
	SampleRate int     `json:"sample_rate"`
	Channels   int     `json:"channels"`
	Duration   float64 `json:"duration"`
}

type VideoMetadataSerializer struct {
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	FrameRate int     `json:"frame_rate"`
	Duration  float64 `json:"duration"`
}

type ImageMetadataSerializer struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type MediaSourcePublicSerializer struct {
	StorageType models.MediaSourceStorageType `json:"storage_type"`
	URL         string                        `json:"url"`
	FileType    models.MediaSourceFileType    `json:"file_type"`
	FormatType  models.MediaSourceFormatType  `json:"format_type"`
	Checksum    string                        `json:"checksum"`
	AudioMeta   *AudioMetadataSerializer      `json:"audio_metadata,omitempty"`
	VideoMeta   *VideoMetadataSerializer      `json:"video_metadata,omitempty"`
	ImageMeta   *ImageMetadataSerializer      `json:"image_metadata,omitempty"`
}

type MediaSourceSerializer struct {
	ID          uint                          `json:"id"`
	StorageType models.MediaSourceStorageType `json:"storage_type"`
	URL         string                        `json:"url"`
	FileType    models.MediaSourceFileType    `json:"file_type"`
	FormatType  models.MediaSourceFormatType  `json:"format_type"`
	Checksum    string                        `json:"checksum"`
	SongHashID  string                        `json:"song_hash_id"`
	AudioMeta   *AudioMetadataSerializer      `json:"audio_metadata,omitempty"`
	VideoMeta   *VideoMetadataSerializer      `json:"video_metadata,omitempty"`
	ImageMeta   *ImageMetadataSerializer      `json:"image_metadata,omitempty"`
}

func SerializeAudioMetadata(m models.AudioMetadata) *AudioMetadataSerializer {
	return &AudioMetadataSerializer{
		SampleRate: m.SampleRate,
		Channels:   m.Channels,
		Duration:   m.Duration,
	}
}

func SerializeVideoMetadata(m models.VideoMetadata) *VideoMetadataSerializer {
	return &VideoMetadataSerializer{
		Width:     m.Width,
		Height:    m.Height,
		FrameRate: m.FrameRate,
		Duration:  m.Duration,
	}
}

func SerializeImageMetadata(m models.ImageMetadata) *ImageMetadataSerializer {
	return &ImageMetadataSerializer{
		Width:  m.Width,
		Height: m.Height,
	}
}

func SerializeMediaSourcePublic(ms models.MediaSource) MediaSourcePublicSerializer {
	var audio *AudioMetadataSerializer
	var video *VideoMetadataSerializer
	var image *ImageMetadataSerializer

	if ms.AudioMetadata != nil {
		audio = &AudioMetadataSerializer{
			SampleRate: ms.AudioMetadata.SampleRate,
			Channels:   ms.AudioMetadata.Channels,
			Duration:   ms.AudioMetadata.Duration,
		}
	}

	if ms.VideoMetadata != nil {
		video = &VideoMetadataSerializer{
			Width:     ms.VideoMetadata.Width,
			Height:    ms.VideoMetadata.Height,
			FrameRate: ms.VideoMetadata.FrameRate,
			Duration:  ms.VideoMetadata.Duration,
		}
	}

	if ms.ImageMetadata != nil {
		image = &ImageMetadataSerializer{
			Width:  ms.ImageMetadata.Width,
			Height: ms.ImageMetadata.Height,
		}
	}

	return MediaSourcePublicSerializer{
		StorageType: ms.StorageType,
		URL:         ms.URL,
		FileType:    ms.FileType,
		FormatType:  ms.FormatType,
		Checksum:    ms.Checksum,
		AudioMeta:   audio,
		VideoMeta:   video,
		ImageMeta:   image,
	}
}

func SerializeMediaSource(ms models.MediaSource) MediaSourceSerializer {
	public := SerializeMediaSourcePublic(ms)
	return MediaSourceSerializer{
		ID:          ms.ID,
		StorageType: public.StorageType,
		URL:         public.URL,
		FileType:    public.FileType,
		FormatType:  public.FormatType,
		Checksum:    public.Checksum,
		SongHashID:  ms.SongHashID,
		AudioMeta:   public.AudioMeta,
		VideoMeta:   public.VideoMeta,
		ImageMeta:   public.ImageMeta,
	}
}
