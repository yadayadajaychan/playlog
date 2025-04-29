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

	"github.com/yadayadajaychan/playlog/database"
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

	songdb, err := database.NewSongDB(db)
	if err != nil {
		panic(err)
	}

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
		addSong(songdb, song)
	}
}

func addSong(songdb *database.SongDB, song songData) {
	songInfo := database.SongInfo{
		Song_id:  song.Song_id,
		Name:     song.Name,
		Artist:   song.Artist,
		Type:     song.Type,
		Bpm:      int(song.Bpm),
		Category: song.Category,
		Version:  song.Version,
		Sort:     song.Sort,
		Charts:   make([]database.ChartInfo, 0),
	}

	err := songdb.AddSong(songInfo)
	if err != nil {
		panic(err)
	}
}
