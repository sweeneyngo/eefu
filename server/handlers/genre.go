package handlers

import (
	"eefu/models"
	"eefu/response"
	"eefu/types"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func GetGenres(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var genres []models.Genre
		if err := db.Find(&genres).Error; err != nil {
			response.InternalServerError(err, "failed to fetch genres").Respond(w, r)
			return
		}
		response.RespondWithJSON(w, http.StatusOK, genres)
	}
}

func CreateGenre(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input types.GenreInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.BadRequest("invalid JSON").Respond(w, r)
			return
		}

		if err := Validate.Struct(input); err != nil {
			response.BadRequest("validation failed").Respond(w, r)
			return
		}

		var existing models.Genre
		if err := db.Where("name = ?", input.Name).First(&existing).Error; err == nil {
			response.BadRequest("genre with this name already exists").Respond(w, r)
			return
		}

		genre := models.Genre{
			Name: input.Name,
		}
		if err := db.Create(&genre).Error; err != nil {
			response.InternalServerError(err, "failed to create genre").Respond(w, r)
			return
		}

		response.RespondWithJSON(w, http.StatusCreated, genre)
	}
}
