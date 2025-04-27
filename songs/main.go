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

// script to store data from songs.json in sqlite3 db
package main

import (
	"os"
	"encoding/json"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type songData struct {
	Song_id		int
	Name		string
	Artist		string
	Type		string
	Bpm		float32
	Category	string
	Version		string
	Sort		string
	Charts		[]struct {
		Difficulty	string
		Level		int
		Internal_level	float32
		Notes_designer	string
		Max_notes	int
	}
}

func main() {
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		panic(err)
	}
	defer db.Close()

	initDB(db)

	file, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	songs := make([]songData, 0, 2048)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&songs)
	if err != nil {
		panic(err)
	}

	for _, song := range songs {
		addSong(db, song)
	}
}

func addSong(db *sql.DB, song songData) {
	_, err := db.Exec(`
	INSERT OR IGNORE INTO songs (
		song_id, name, artist, type,
		bpm, category, version, sort) VALUES (
		?, ?, ?, ?,
		?, ?, ?, ?
	);`,
		song.Song_id, song.Name, song.Artist, song.Type,
		int(song.Bpm), song.Category, song.Version, song.Sort)
	if err != nil {
		panic(err)
	}
}

func initDB(db *sql.DB) {
	_, err := db.Exec(`
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
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS basic (
		song_id        INTEGER PRIMARY KEY NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS advanced (
		song_id        INTEGER PRIMARY KEY NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS expert (
		song_id        INTEGER PRIMARY KEY NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS master (
		song_id        INTEGER PRIMARY KEY NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS remaster (
		song_id        INTEGER PRIMARY KEY NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER
	);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS utage (
		song_id        INTEGER PRIMARY KEY NOT NULL,
		level          INTEGER,
		internal_level REAL,
		notes_designer TEXT,
		max_notes      INTEGER
	);`)
	if err != nil {
		panic(err)
	}
}
