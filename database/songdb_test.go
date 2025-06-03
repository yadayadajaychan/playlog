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
package database_test

import (
	"testing"
	"os"
	"reflect"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yadayadajaychan/playlog/database"
)

func TestAddAndGetSong(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "songdb-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	db, err := sql.Open("sqlite3", tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	songdb, err := database.NewSongDB(db)
	if err != nil {
		t.Fatal(err)
	}

	song1 := database.SongInfo{
		SongId:   1,
		Name:     "Sweet Home Alabama",
		Artist:   "Lynyrd Skynyrd",
		Type:     "std",
		Bpm:      98,
		Category: "POPS＆アニメ",
		Version:  "maimai",
		Sort:     "100000",
		Charts:   make([]database.ChartInfo, 0, 5),
	}

	song1.Charts = append(song1.Charts, database.ChartInfo{
		Difficulty:    database.Basic,
		Level:         6,
		InternalLevel: 69,
		NotesDesigner: "your mom",
		MaxNotes: 96,
	})

	song1.Charts = append(song1.Charts, database.ChartInfo{
		Difficulty:    database.Advanced,
		Level:         8,
		InternalLevel: 88,
		NotesDesigner: "mein kampf",
		MaxNotes: 1488,
	})

	song1.Charts = append(song1.Charts, database.ChartInfo{
		Difficulty:    database.Expert,
		Level:         10,
		InternalLevel: 102,
		NotesDesigner: "",
		MaxNotes: 1600,
	})

	err = songdb.AddSong(song1)
	if err != nil {
		t.Fatal(err)
	}

	song1g, err := songdb.GetSong(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = songdb.GetSong(2)
	if _, ok := err.(*database.SongNotFoundError); ok {
		t.Log("correctly returned SongNotFoundError:", err)
	} else if err != nil {
		t.Error("returned non-nil error, but not SongNotFoundError:", err)
	} else {
		t.Error("expected SongNotFoundError for non-existant songId")
	}

	if !reflect.DeepEqual(song1, song1g) {
		t.Log(song1)
		t.Log(song1g)
		t.Error("song1 not equal")
	}
}

func TestGetSongsByName(t *testing.T) {
	db, err := sql.Open("sqlite3", "../songs.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	songdb, err := database.NewSongDB(db)
	if err != nil {
		t.Fatal(err)
	}

	song1, err := songdb.GetSongsByName("終焉逃避行")
	if err != nil {
		t.Fatal("error retrieving song1")
	}
	if song1[0].SongId != 11441 {
		t.Error("song1: incorrect song retrieved")
	}

	song2, err := songdb.GetSongsByName("Scatman (Ski Ba Bop Ba Dop Bop)")
	if err != nil {
		t.Fatal("error retrieving song2")
	}
	if song2[0].SongId != 502 {
		t.Error("song2: incorrect song retrieved")
	}

	song3, err := songdb.GetSongsByName("天ノ弱")
	if err != nil {
		t.Fatal("error retrieving song3")
	}
	if len(song3) != 2 {
		t.Fatal("song3: incorrect no. of songs returned")
	}
	if song3[0].SongId != 188 && song3[0].SongId != 10188 {
		t.Error("song3: incorrect song retrieved")
	}
	if song3[1].SongId != 188 && song3[1].SongId != 10188 {
		t.Error("song3: incorrect song retrieved")
	}

	song4, err := songdb.GetSongsByName("i went to ur mom's house")
	if err != nil {
		t.Fatal("error retrieving song4:", err)
	}
	if len(song4) != 0 {
		t.Error("song4: incorrect no. of songs returned:", len(song4))
	}
}
