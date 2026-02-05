package handlers

import (
	"eefu/models"
	"eefu/response"
	"eefu/types"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func GetSingers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var singers []models.Singer
		if err := db.Preload("Aliases").Find(&singers).Error; err != nil {
			response.InternalServerError(err, "failed to fetch singers").Respond(w, r)
			return
		}
		response.RespondWithJSON(w, http.StatusOK, singers)
	}
}

func CreateSinger(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input types.SingerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.BadRequest("invalid JSON").Respond(w, r)
			return
		}

		if err := Validate.Struct(input); err != nil {
			response.BadRequest("validation failed").Respond(w, r)
			return
		}

		var existing models.Singer
		if err := db.Where("name = ?", input.Name).First(&existing).Error; err == nil {
			response.BadRequest("singer with this name already exists").Respond(w, r)
			return
		}

		var singer models.Singer

		err := db.Transaction(func(tx *gorm.DB) error {
			singer = models.Singer{Name: input.Name}
			if err := tx.Create(&singer).Error; err != nil {
				return err
			}

			for _, a := range input.Aliases {
				alias := models.SingerAlias{
					SingerID: singer.ID,
					Name:     a.Name,
					Language: a.Language,
				}
				if err := tx.Create(&alias).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			response.InternalServerError(err, "failed to create singer").Respond(w, r)
			return
		}

		db.Preload("Aliases").First(&singer, singer.ID)
		response.RespondWithJSON(w, http.StatusCreated, singer)
	}
}
