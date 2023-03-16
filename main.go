package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
)

type Song struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
	Singer   string  `json:"singer"`
}

func main() {
	// Connect to database
	db, err := sql.Open("mysql", "user:@tcp(127.0.0.1:3306)/quiz")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m := martini.Classic()

	// Get all songs
	m.Get("/songs", func() (int, string) {
		var songs []Song
		rows, err := db.Query("SELECT id, title, duration, singer FROM songs")
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}
		defer rows.Close()

		for rows.Next() {
			var song Song
			err := rows.Scan(&song.ID, &song.Title, &song.Duration, &song.Singer)
			if err != nil {
				log.Println(err)
				return http.StatusInternalServerError, "Internal server error"
			}
			songs = append(songs, song)
		}

		if len(songs) == 0 {
			return http.StatusNotFound, "No songs found"
		}

		return http.StatusOK, fmt.Sprintf("%+v", songs)
	})

	// Get a song by ID
	m.Get("/songs/:id", func(params martini.Params) (int, string) {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			return http.StatusBadRequest, "Invalid song ID"
		}

		var song Song
		err = db.QueryRow("SELECT id, title, duration, singer FROM songs WHERE id = ?", id).Scan(&song.ID, &song.Title, &song.Duration, &song.Singer)
		if err != nil {
			if err == sql.ErrNoRows {
				return http.StatusNotFound, "Song not found"
			}
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		return http.StatusOK, fmt.Sprintf("%+v", song)
	})

	// Add a new song
	m.Post("/songs", func(r *http.Request) (int, string) {
		title := r.FormValue("title")
		if title == "" {
			return http.StatusBadRequest, "Song title is required"
		}

		duration, err := strconv.ParseFloat(r.FormValue("duration"), 64)
		if err != nil {
			return http.StatusBadRequest, "Invalid song duration"
		}

		singer := r.FormValue("singer")
		if singer == "" {
			return http.StatusBadRequest, "Song singer is required"
		}

		res, err := db.Exec("INSERT INTO songs (title, duration, singer) VALUES (?, ?, ?)", title, duration, singer)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		return http.StatusCreated, fmt.Sprintf("Song added with ID %d", id)
	})

	// Update a song by ID
	m.Put("/songs/:id", func(params martini.Params, r *http.Request) (int, string) {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			return http.StatusBadRequest, "Invalid song ID"
		}
		title := r.FormValue("title")
		if title == "" {
			return http.StatusBadRequest, "Song title is required"
		}

		duration, err := strconv.ParseFloat(r.FormValue("duration"), 64)
		if err != nil {
			return http.StatusBadRequest, "Invalid song duration"
		}

		singer := r.FormValue("singer")
		if singer == "" {
			return http.StatusBadRequest, "Song singer is required"
		}

		res, err := db.Exec("UPDATE songs SET title = ?, duration = ?, singer = ? WHERE id = ?", title, duration, singer, id)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		if rowsAffected == 0 {
			return http.StatusNotFound, "Song not found"
		}

		return http.StatusOK, fmt.Sprintf("Song with ID %d updated", id)
	})

	// Delete a song by ID
	m.Delete("/songs/:id", func(params martini.Params) (int, string) {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			return http.StatusBadRequest, "Invalid song ID"
		}

		res, err := db.Exec("DELETE FROM songs WHERE id = ?", id)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		if rowsAffected == 0 {
			return http.StatusNotFound, "Song not found"
		}

		return http.StatusOK, fmt.Sprintf("Song with ID %d deleted", id)
	})

	m.Run()
}
