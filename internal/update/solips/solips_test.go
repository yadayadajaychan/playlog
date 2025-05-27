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

package solips

import (
	"os"
	"testing"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yadayadajaychan/playlog/database"
)

// TestGetPlaylog tests that the returned playlog has 100 entries
// with unique PlaylogApiIds and UserPlayDates
func TestGetPlaylog(t *testing.T) {
	accessCode := os.Getenv("PLAYLOG_ACCESS_CODE")
	if accessCode == "" {
		t.Fatal("missing 'PLAYLOG_ACCESS_CODE' environment variable")
	}

	db, err := sql.Open("sqlite3", "../../../songs.db")
	if err != nil {
	        t.Fatal(err)
	}
	defer db.Close()

	songdb, err := database.NewSongDB(db)
	if err != nil {
		t.Fatal(err)
	}

	playlog, err := getPlaylog(accessCode)
	if err != nil {
		t.Fatal(err)
	}

	err = validatePlaylog(playlog)
	if err != nil {
		t.Fatal(err)
	}

	playlogDetail, err := getPlaylogDetail(accessCode, playlog.Playlog[0].PlaylogApiId)
	if err != nil {
		t.Fatal(err)
	}

	err = validatePlaylogDetail(playlogDetail, songdb)
	if err != nil {
		t.Fatal(err)
	}
}
