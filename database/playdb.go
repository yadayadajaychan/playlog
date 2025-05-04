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
)

type PlayDB struct {
	db *sql.DB
}

func NewPlayDB(db *sql.DB) (*PlayDB, error) {
	playdb := &PlayDB{db: db}
	err := playdb.initDB()
	return playdb, err
}

func (playdb *PlayDB) initDB() error {
	tx, err := playdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS plays (
		user_play_date INTEGER PRIMARY KEY NOT NULL,
		song_id        INTEGER NOT NULL,
		difficulty     INTEGER NOT NULL,

		score            INTEGER,
		dx_score         INTEGER,
		combo_status     INTEGER,
		sync_status      INTEGER,
		is_clear         INTEGER,
		is_new_record    INTEGER,
		is_dx_new_record INTEGER,
		track            INTEGER,
		matching_users   TEXT,

		max_combo   INTEGER,
		total_combo INTEGER,
		max_sync    INTEGER,
		total_sync  INTEGER,

		fast_count    INTEGER,
		late_count    INTEGER,
		before_rating INTEGER,
		after_rating  INTEGER,

		tap_critical_perfect INTEGER,
		tap_perfect          INTEGER,
		tap_great            INTEGER,
		tap_good             INTEGER,
		tap_miss             INTEGER,

		hold_critical_perfect INTEGER,
		hold_perfect          INTEGER,
		hold_great            INTEGER,
		hold_good             INTEGER,
		hold_miss             INTEGER,

		slide_critical_perfect INTEGER,
		slide_perfect          INTEGER,
		slide_great            INTEGER,
		slide_good             INTEGER,
		slide_miss             INTEGER,

		touch_critical_perfect INTEGER,
		touch_perfect          INTEGER,
		touch_great            INTEGER,
		touch_good             INTEGER,
		touch_miss             INTEGER,

		break_critical_perfect INTEGER,
		break_perfect          INTEGER,
		break_great            INTEGER,
		break_good             INTEGER,
		break_miss             INTEGER,

		total_critical_perfect INTEGER,
		total_perfect          INTEGER,
		total_great            INTEGER,
		total_good             INTEGER,
		total_miss             INTEGER
	);`)
	if err != nil {
		return err
	}

	return tx.Commit()
}
