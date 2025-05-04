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
	if err != nil {
		t.Log("correctly returned non-nil error:", err)
	} else {
		t.Error("expected non-nil error for non-existant songId")
	}

	if !reflect.DeepEqual(song1, song1g) {
		t.Log(song1)
		t.Log(song1g)
		t.Error("song1 not equal")
	}
}
