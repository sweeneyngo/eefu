package types

import "eefu/models"

type MediaMinimal struct {
	FileType models.MediaSourceFileType `json:"file_type"`
	URL      string                     `json:"url"`
}
