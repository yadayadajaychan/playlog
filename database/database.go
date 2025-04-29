// Copyright (C) 2025 Ethan Cheng <ethan@nijika.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package database handles the playlog and song database
package database

import (
	"database/sql"
	"errors"
)

type SongInfo struct {
	Song_id		int
	Name		string
	Artist		string
	Type		string
	Bpm		int
	Category	string
	Version		string
	Sort		string
	Charts		[]ChartInfo
}

type ChartInfo struct {
	Difficulty	Difficulty
	Level		int
	Internal_level	float32
	Notes_designer	string
	Max_notes	int
}

type Difficulty int

const (
	Basic Difficulty = iota
	Advanced
	Expert
	Master
	ReMaster
	Utage
)

type SongDB struct {
	db	*sql.DB
}

// NewSongDB creates a SongDB object and initializes the database
func NewSongDB(db *sql.DB) (*SongDB, error) {
	songdb := &SongDB{db: db}
	err := songdb.initDB()
	return songdb, err
}

func (songdb *SongDB) initDB() error {
	tx, err := songdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS songs (
		song_id  INTEGER PRIMARY KEY NOT NULL,
		name     TEXT,
		artist   TEXT,
		type     TEXT,
		bpm      INTEGER,
		category TEXT,
		version  TEXT,
		sort     TEXT
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS charts (
		song_id        INTEGER NOT NULL,
		difficulty     INTEGER NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER,
		PRIMARY KEY (song_id, difficulty)
	);`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// AddSong adds a song to the song db, ignoring if the song already exists
func (songdb *SongDB) AddSong(song SongInfo) error {
	tx, err := songdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	INSERT OR IGNORE INTO songs (
		song_id, name, artist, type,
		bpm, category, version, sort) VALUES (
		?, ?, ?, ?,
		?, ?, ?, ?
	);`,
		song.Song_id, song.Name, song.Artist, song.Type,
		song.Bpm, song.Category, song.Version, song.Sort)
	if err != nil {
		return err
	}

	for _, chart := range song.Charts {
		_, err = tx.Exec(`
		INSERT OR IGNORE INTO charts (
			song_id, difficulty, level,
			internal_level, notes_designer, max_notes) VALUES (
			?, ?, ?,
			?, ?, ?
		);`,
			song.Song_id, chart.Difficulty, chart.Level,
			chart.Internal_level, chart.Notes_designer,
			chart.Max_notes)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (songdb *SongDB) GetSong(songId int) (SongInfo, error) {
	song := SongInfo{}

	rows, err := songdb.db.Query(`
	SELECT song_id, name, artist, type,
		bpm, category, version, sort FROM songs WHERE song_id=?`, songId)
	if err != nil {
		return song, err
	}

	if rows.Next() {
		err = rows.Scan(&song.Song_id, &song.Name, &song.Artist, &song.Type,
				&song.Bpm, &song.Category, &song.Version, &song.Sort)
		if err != nil {
			return song, err
		}
	} else {
		return song, errors.New("no song found")
	}
	// TODO: custom error type

	rows, err = songdb.db.Query(`
	SELECT difficulty, level, internal_level,
		notes_designer, max_notes FROM charts WHERE song_id=?`, songId)
	if err != nil {
		return song, err
	}

	for rows.Next() {
		chart := ChartInfo{}
		err = rows.Scan(&chart.Difficulty, &chart.Level, &chart.Internal_level,
				&chart.Notes_designer, &chart.Max_notes)
		if err != nil {
			return song, err
		}

		song.Charts = append(song.Charts, chart)
	}

	if song.Charts == nil {
		return song, errors.New("no chart found")
	}

	return song, nil
}
