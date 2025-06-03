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
	"fmt"
)

type SongDB struct {
	db *sql.DB
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
		internal_level INTEGER,
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
		song.SongId, song.Name, song.Artist, song.Type,
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
			song.SongId, chart.Difficulty, chart.Level,
			chart.InternalLevel, chart.NotesDesigner,
			chart.MaxNotes)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// takes row with one item
// caller's responsibility to call Close() on rows
func (songdb *SongDB) rowsToSongInfos(rows *sql.Rows) ([]SongInfo, error) {
	songs := make([]SongInfo, 0, 1)

	for rows.Next() {
		song := SongInfo{}
		err := rows.Scan(&song.SongId, &song.Name, &song.Artist, &song.Type,
				&song.Bpm, &song.Category, &song.Version, &song.Sort)
		if err != nil {
			return songs, err
		}

		chartRows, err := songdb.db.Query(`
			SELECT difficulty, level, internal_level,
			notes_designer, max_notes FROM charts WHERE song_id=?`, song.SongId)
		if err != nil {
			return songs, err
		}
		defer chartRows.Close()

		for chartRows.Next() {
			chart := ChartInfo{}
			err = chartRows.Scan(&chart.Difficulty, &chart.Level, &chart.InternalLevel,
					&chart.NotesDesigner, &chart.MaxNotes)
			if err != nil {
				return songs, err
			}

			song.Charts = append(song.Charts, chart)
		}

		if len(song.Charts) == 0 {
			return songs, errors.New("no chart found")
		}

		songs = append(songs, song)
	}

	return songs, nil
}

// GetSong gets a song from the database using the songId
func (songdb *SongDB) GetSong(songId int) (SongInfo, error) {
	rows, err := songdb.db.Query(`
		SELECT * FROM songs WHERE song_id=?`, songId)
	if err != nil {
		return SongInfo{}, err
	}
	defer rows.Close()

	songs, err := songdb.rowsToSongInfos(rows)
	if err != nil {
		return SongInfo{}, err
	}

	if len(songs) <= 0 {
		return SongInfo{}, &SongNotFoundError{SongId: songId}
	}

	return songs[0], nil
}

// GetSongsByName returns songs from the database using 'name'
// Can return both the std and dx versions
func (songdb *SongDB) GetSongsByName(name string) ([]SongInfo, error) {
	rows, err := songdb.db.Query(`
		SELECT * FROM songs WHERE name=?`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs, err := songdb.rowsToSongInfos(rows)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

// GetSongsByVersion returns songs from the database with the same version
func (songdb *SongDB) GetSongsByVersion(version string) ([]SongInfo, error) {
	rows, err := songdb.db.Query(`
		SELECT * FROM songs WHERE version=? COLLATE NOCASE`, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs, err := songdb.rowsToSongInfos(rows)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

type SongNotFoundError struct {
	SongId	int
	Name    string
}

func (e *SongNotFoundError) Error() string {
	return fmt.Sprintf("Song with id %d / name '%s' not found in database", e.SongId, e.Name)
}
