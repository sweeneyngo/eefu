package serializers

import (
	"eefu/models"
)

type TagPublicSerializer struct {
	Name        string         `json:"name"`
	Type        models.TagType `json:"type"`
	Description string         `json:"description"`
}

type TagSerializer struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Type        models.TagType `json:"type"`
	Description string         `json:"description"`
}

func SerializeTagPublic(tag models.Tag) TagPublicSerializer {
	return TagPublicSerializer{
		Name:        tag.Name,
		Type:        tag.Type,
		Description: tag.Description,
	}
}

func SerializeTag(tag models.Tag) TagSerializer {
	return TagSerializer{
		ID:          tag.ID,
		Name:        tag.Name,
		Type:        tag.Type,
		Description: tag.Description,
	}
}
