package handlers

import (
	"eefu/models"
	"eefu/response"
	"eefu/types"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func GetTags(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tags []models.Tag
		if err := db.Find(&tags).Error; err != nil {
			response.InternalServerError(err, "failed to fetch tags").Respond(w, r)
			return
		}
		response.RespondWithJSON(w, http.StatusOK, tags)
	}
}

func CreateTag(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input types.TagInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.BadRequest("invalid JSON").Respond(w, r)
			return
		}

		if err := Validate.Struct(input); err != nil {
			response.BadRequest("validation failed").Respond(w, r)
			return
		}

		var existing models.Tag
		if err := db.Where("name = ? AND type = ?", input.Name, input.Type).First(&existing).Error; err == nil {
			response.BadRequest("tag with this name and type already exists").Respond(w, r)
			return
		}

		tag := models.Tag{
			Name: input.Name,
			Type: models.TagType(input.Type),
		}
		if input.Description != nil {
			tag.Description = *input.Description
		}

		if err := db.Create(&tag).Error; err != nil {
			response.InternalServerError(err, "failed to create tag").Respond(w, r)
			return
		}

		response.RespondWithJSON(w, http.StatusCreated, tag)
	}
}
