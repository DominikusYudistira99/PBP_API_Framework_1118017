package Model

import "database/sql"

type Song struct {
	SongID       int64   `json:"song_id"`
	SongTitle    string  `json:"song_title"`
	SongDuration float64 `json:"song_duration"`
	SongSinger   string  `json:"song_singer"`
}

type SongModel struct {
	DB *sql.DB
}

func (m *SongModel) Create(song *Song) error {
	query := "INSERT INTO songs (song_title, song_duration, song_singer) VALUES (?, ?, ?)"
	result, err := m.DB.Exec(query, song.SongTitle, song.SongDuration, song.SongSinger)
	if err != nil {
		return err
	}
	song.SongID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (m *SongModel) Get(id int64) (*Song, error) {
	query := "SELECT * FROM songs WHERE song_id = ?"
	row := m.DB.QueryRow(query, id)
	song := &Song{}
	err := row.Scan(&song.SongID, &song.SongTitle, &song.SongDuration, &song.SongSinger)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (m *SongModel) GetAll() ([]*Song, error) {
	query := "SELECT * FROM songs"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := []*Song{}
	for rows.Next() {
		song := &Song{}
		err := rows.Scan(&song.SongID, &song.SongTitle, &song.SongDuration, &song.SongSinger)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

func (m *SongModel) Update(song *Song) error {
	query := "UPDATE songs SET song_title = ?, song_duration = ?, song_singer = ? WHERE song_id = ?"
	_, err := m.DB.Exec(query, song.SongTitle, song.SongDuration, song.SongSinger, song.SongID)
	if err != nil {
		return err
	}
	return nil
}

func (m *SongModel) Delete(id int64) error {
	query := "DELETE FROM songs WHERE song_id = ?"
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
