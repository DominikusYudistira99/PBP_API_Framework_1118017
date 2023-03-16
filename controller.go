package controllers

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/go-martini/martini"
)

type SongController struct {
    SongModel *model.SongModel
}

func (sc *SongController) CreateSong(w http.ResponseWriter, r *http.Request) {
    var song model.Song
    err := json.NewDecoder(r.Body).Decode(&song)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    err = sc.SongModel.Create(&song)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

