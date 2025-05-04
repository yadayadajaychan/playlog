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
	"encoding/json"
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

func validatePlay(play PlayInfo) error {
	// TODO
	return nil
}

func (playdb *PlayDB) AddPlay(play PlayInfo) error {
	err := validatePlay(play)
	if err != nil {
		return err
	}

	matchingUsersJSON, err := json.Marshal(play.MatchingUsers)
	if err != nil {
		return err
	}

	tx, err := playdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	INSERT OR IGNORE INTO plays (
		user_play_date, song_id, difficulty,

		score, dx_score, combo_status, sync_status,
		is_clear, is_new_record, is_dx_new_record,
		track, matching_users,

		max_combo, total_combo, max_sync, total_sync,

		fast_count, late_count, before_rating, after_rating,

		tap_critical_perfect, tap_perfect, tap_great,
		tap_good, tap_miss,

		hold_critical_perfect, hold_perfect, hold_great,
		hold_good, hold_miss,

		slide_critical_perfect, slide_perfect, slide_great,
		slide_good, slide_miss,

		touch_critical_perfect, touch_perfect, touch_great,
		touch_good, touch_miss,

		break_critical_perfect, break_perfect, break_great,
		break_good, break_miss,

		total_critical_perfect, total_perfect, total_great,
		total_good, total_miss
	) VALUES (
		?, ?, ?,

		?, ?, ?, ?,
		?, ?, ?,
		?, ?,

		?, ?, ?, ?,

		?, ?, ?, ?,

		?, ?, ?,
		?, ?,

		?, ?, ?,
		?, ?,

		?, ?, ?,
		?, ?,

		?, ?, ?,
		?, ?,

		?, ?, ?,
		?, ?,

		?, ?, ?,
		?, ?
	);`,
		play.UserPlayDate, play.SongId, play.Difficulty,

		play.Score, play.DxScore, play.ComboStatus, play.SyncStatus,
		play.IsClear, play.IsNewRecord, play.IsDxNewRecord,
		play.Track, matchingUsersJSON,

		play.MaxCombo, play.TotalCombo, play.MaxSync, play.TotalSync,

		play.FastCount, play.LateCount, play.BeforeRating, play.AfterRating,

		play.TapCriticalPerfect, play.TapPerfect, play.TapGreat,
		play.TapGood, play.TapMiss,

		play.HoldCriticalPerfect, play.HoldPerfect, play.HoldGreat,
		play.HoldGood, play.HoldMiss,

		play.SlideCriticalPerfect, play.SlidePerfect, play.SlideGreat,
		play.SlideGood, play.SlideMiss,

		play.TouchCriticalPerfect, play.TouchPerfect, play.TouchGreat,
		play.TouchGood, play.TouchMiss,

		play.BreakCriticalPerfect, play.BreakPerfect, play.BreakGreat,
		play.BreakGood, play.BreakMiss,

		play.TotalCriticalPerfect, play.TotalPerfect, play.TotalGreat,
		play.TotalGood, play.TotalMiss)

	if err != nil {
		return err
	}

	return tx.Commit()
}
