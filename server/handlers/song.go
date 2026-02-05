package handlers

import (
	"eefu/middleware"
	"eefu/models"
	"eefu/response"
	"eefu/serializers"
	"eefu/services"
	"eefu/types"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func GetSongs(svc *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var songs []models.Song
		songs, err := svc.GetAll(r.Context())
		if err != nil {
			response.InternalServerError(err, "failed to fetch songs").Respond(w, r)
			return
		}
		response.RespondWithJSON(w, http.StatusOK, songs)
	}
}

func GetSongVersionsByGroup(svc *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songGroupHashID := chi.URLParam(r, "song_group_hash_id")
		if songGroupHashID == "" {
			response.BadRequest("song_group_hash_id param is required").Respond(w, r)
			return
		}

		versions, err := svc.GetVersionsByGroup(r.Context(), songGroupHashID)
		if err != nil {
			response.InternalServerError(err, "failed to fetch song versions").Respond(w, r)
			return
		}

		if len(versions) == 0 {
			response.NotFound("no songs found for this group").Respond(w, r)
			return
		}

		isAdmin, _ := r.Context().Value(middleware.IsAdminKey).(bool)
		serialized := make([]any, len(versions))
		for i, s := range versions {
			if isAdmin {
				serialized[i] = serializers.SerializeSong(s)
			} else {
				serialized[i] = serializers.SerializeSongPublic(s)
			}
		}

		response.RespondWithJSON(w, http.StatusOK, serialized)
	}
}

func GetSongMedia(svc *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashID := chi.URLParam(r, "hash_id")
		if hashID == "" {
			response.BadRequest("hash_id param is required").Respond(w, r)
			return
		}

		mediaList, err := svc.GetMedia(r.Context(), hashID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				response.BadRequest("song not found").Respond(w, r)
				return
			}

			response.InternalServerError(err, "failed to fetch song media").Respond(w, r)
			return
		}

		response.RespondWithJSON(w, http.StatusOK, mediaList)
	}
}

func CreateSong(svc *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input types.SongInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.BadRequest("invalid JSON").Respond(w, r)
			return
		}

		if err := Validate.Struct(input); err != nil {
			response.BadRequest("validation failed").Respond(w, r)
			return
		}

		song, err := svc.CreateSong(r.Context(), input)
		if err != nil {
			response.InternalServerError(err, "failed to create song").Respond(w, r)
			return
		}
		response.RespondWithJSON(w, http.StatusOK, song)
	}
}

func CreateSongVersion(svc *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashID := chi.URLParam(r, "hash_id")

		var input types.SongVersionInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.BadRequest("invalid JSON").Respond(w, r)
			return
		}

		if err := Validate.Struct(input); err != nil {
			response.BadRequest("validation failed").Respond(w, r)
			return
		}

		log.Print(hashID)
		song, err := svc.CreateSongVersion(r.Context(), hashID, input)
		if err != nil {
			response.InternalServerError(err, "failed to create song").Respond(w, r)
			return
		}
		response.RespondWithJSON(w, http.StatusOK, song)
	}
}

func UploadMedia(svc *services.MediaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hashID := chi.URLParam(r, "hash_id")

		if err := r.ParseMultipartForm(100 << 20); err != nil {
			response.BadRequest("invalid multipart form").Respond(w, r)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			response.BadRequest("file not found").Respond(w, r)
			return
		}
		defer file.Close()

		fileType := models.MediaSourceFileType(r.FormValue("file_type"))
		storageType := models.MediaSourceStorageType(r.FormValue("storage_type"))

		media, err := svc.UploadMedia(ctx, hashID, file, header.Filename, fileType, storageType)
		if err != nil {
			response.InternalServerError(err, "failed to upload media").Respond(w, r)
			return
		}

		response.RespondWithJSON(w, http.StatusCreated, media)
	}
}
