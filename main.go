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

package main

import (
	"time"
	"os"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yadayadajaychan/playlog/update"
	"github.com/yadayadajaychan/playlog/database"
	_ "github.com/pborman/getopt/v2"
)

func main() {
	accessCode := os.Getenv("PLAYLOG_ACCESS_CODE")
	if accessCode == "" {
		log.Fatal("missing 'PLAYLOG_ACCESS_CODE' environment variable")
	}

	db, err := sql.Open("sqlite3", "tmp.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	playdb, err := database.NewPlayDB(db)
	if err != nil {
		panic(err)
	}

	update.Update(playdb, accessCode, 1 * time.Second)
}
