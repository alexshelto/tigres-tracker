package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexshelto/tigres-tracker/api/models"
)

func AddSongPlay(w http.ResponseWriter, r *http.Request) {
	var songPlay models.SongPlayRequest

	err := json.NewDecoder(r.Body).Decode(&songPlay)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	line := fmt.Sprintf("received: %+v", songPlay)
	w.Write([]byte(line))
}
