package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/alexshelto/tigres-tracker/api/database"
	"github.com/alexshelto/tigres-tracker/api/dto"
	"github.com/alexshelto/tigres-tracker/api/handlers"
	"github.com/alexshelto/tigres-tracker/api/models"
	"github.com/alexshelto/tigres-tracker/api/repository"
	"github.com/alexshelto/tigres-tracker/api/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetDBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", Db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var Db *gorm.DB

func main() {
	Db = database.ConnectDb()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(SetDBMiddleware)

	userRepo := &repository.UserRepository{}
	songRepo := &repository.SongRepository{}
	songPlayRepo := &repository.SongPlayRepository{}

	songPlayService := service.NewSongPlayService(userRepo, songRepo, songPlayRepo)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.Index)
	})

	r.Post("/song", func(w http.ResponseWriter, r *http.Request) {
		var songPlay models.SongPlayRequest

		err := json.NewDecoder(r.Body).Decode(&songPlay)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = songPlayService.AddSongPlay(r.Context(), songPlay)

		if err != nil {
			http.Error(w, fmt.Sprintf("Something went wrong %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("success"))
	})

	r.Get("/song/top", func(w http.ResponseWriter, r *http.Request) {
		guildID := r.URL.Query().Get("guild_id")
		limitStr := r.URL.Query().Get("limit")

		if guildID == "" {
			http.Error(w, "guild_id is required", http.StatusBadRequest)
			return
		}

		limit := 10 //Default Limit
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				http.Error(w, "invalid limit value", http.StatusBadRequest)
				return
			}
		}
		topSongs, err := songPlayService.GetTopSongsInGuild(r.Context(), guildID, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(topSongs)
	})

	r.Get("/song/count", func(w http.ResponseWriter, r *http.Request) {
		guildID := r.URL.Query().Get("guild_id")

		if guildID == "" {
			http.Error(w, "guild_id is required", http.StatusBadRequest)
			return
		}

		total, err := songPlayService.GetTotalSongPlaysInGuild(r.Context(), guildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalPlays := dto.TotalSongPlayDTO{TotalPlays: total}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(totalPlays)
	})

	r.Get("/user/{discordID}/song/top", func(w http.ResponseWriter, r *http.Request) {
		discordID := chi.URLParam(r, "discordID")

		if discordID == "" {
			http.Error(w, "discordID is required", http.StatusBadRequest)
			return
		}

		// Extract query parameters
		guildID := r.URL.Query().Get("guild_id")
		if guildID == "" {
			http.Error(w, "guild_id is required", http.StatusBadRequest)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 10 // Default limit
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				http.Error(w, "limit must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		topSongs, err := songPlayService.GetTopSongsByUserInGuild(r.Context(), discordID, guildID, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(topSongs)

	})

	r.Get("/user/{discordID}/song/count", func(w http.ResponseWriter, r *http.Request) {
		discordID := chi.URLParam(r, "discordID")

		if discordID == "" {
			http.Error(w, "discordID is required", http.StatusBadRequest)
			return
		}

		// Extract query parameters
		guildID := r.URL.Query().Get("guild_id")
		if guildID == "" {
			http.Error(w, "guild_id is required", http.StatusBadRequest)
			return
		}

		total, err := songPlayService.GetTotalUserSongPlaysInGuild(r.Context(), discordID, guildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalPlays := dto.TotalSongPlayDTO{TotalPlays: total}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(totalPlays)
	})

	http.ListenAndServe(":3000", r)
}
