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
	"errors"
	"fmt"

	"github.com/yadayadajaychan/playlog/util"
)

type PlayDB struct {
	db *sql.DB
	ascStmt  *sql.Stmt // used to query playlog entries by ascending order
	descStmt *sql.Stmt // used to query playlog entries by descending order
}

func NewPlayDB(db *sql.DB) (*PlayDB, error) {
	playdb := &PlayDB{db: db}
	err := playdb.initDB()
	return playdb, err
}

// called by main init function
func (playdb *PlayDB) initDB_v1_2() error {
	tx, err := playdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS dx_rating_gen_3 (
		user_play_date INTEGER PRIMARY KEY NOT NULL,
		internal_level INTEGER,
		rating         INTEGER,
		version        TEXT
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT OR IGNORE INTO version (major, minor) VALUES (1, 2);`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// PopulateDxRatingGen3 populates the dx_rating_gen_3 table in the db.
// It is usually unnecessary to call this function yourself since it is
// automatically called when trying to access rating info.
func (playdb *PlayDB) PopulateDxRatingGen3(songdb *SongDB) error {
	for {
		rows, err := playdb.db.Query(`
			SELECT * FROM plays WHERE user_play_date NOT IN
			(SELECT user_play_date FROM dx_rating_gen_3) LIMIT 1;`)
		if err != nil {
			return err
		}
		defer rows.Close()

		plays, err := rowsToPlayInfos(rows)
		if err != nil {
			return err
		}
		if len(plays) == 0 {
			break
		}

		play := plays[0]
		song, err := songdb.GetSong(play.SongId)
		if err != nil {
			return err
		}

		var chart *ChartInfo
		for _, c := range song.Charts {
			if c.Difficulty == play.Difficulty {
				chart = &c
			}
		}
		if chart == nil {
			return errors.New("no chart found")
		}

		rating := util.ScoreAndInternalLevelToDxRatingGen3(play.Score, chart.InternalLevel)

		_, err = playdb.db.Exec(`
			INSERT OR IGNORE INTO dx_rating_gen_3
			(user_play_date, internal_level, rating, version)
			VALUES (?, ?, ?, ?);`,
			play.UserPlayDate, chart.InternalLevel,
			rating, song.Version)
		if err != nil {
			return err
		}
	}

	return nil
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

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS version (
		major INTEGER NOT NULL,
		minor INTEGER NOT NULL,
		PRIMARY KEY (major, minor)
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT OR IGNORE INTO version (major, minor) VALUES (1, 0);`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	playdb.ascStmt, err = playdb.db.Prepare(
		`SELECT * FROM plays ORDER BY user_play_date ASC
		LIMIT ? OFFSET ?`)
	if err != nil {
		return err
	}

	playdb.descStmt, err = playdb.db.Prepare(
		`SELECT * FROM plays ORDER BY user_play_date DESC
		LIMIT ? OFFSET ?`)
	if err != nil {
		return err
	}

	err = playdb.initDB_v1_2()
	if err != nil {
		return err
	}

	return nil
}

// GetVersion gets the major and minor version of the play database
// corresponding to the program version that last changed the schema
func (playdb *PlayDB) GetVersion() (major, minor int, err error) {
	rows, err := playdb.db.Query(`
	SELECT * FROM version ORDER BY major DESC, minor DESC LIMIT 1;`)
	if err != nil {
		return major, minor, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&major, &minor)
		if err != nil {
			return major, minor, err
		}
	} else {
		if err = rows.Err(); err != nil {
			return major, minor, err
		} else {
			return 0, 0, nil
		}
	}

	return major, minor, nil
}

func (playdb *PlayDB) GetDxRatingGen3(date int64, songdb *SongDB) (int, error) {
	var rating int

	err := playdb.PopulateDxRatingGen3(songdb)
	if err != nil {
		return rating, err
	}

	rows, err := playdb.db.Query(`
	SELECT rating FROM dx_rating_gen_3 WHERE user_play_date=?`, date)
	if err != nil {
		return rating, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&rating)
		if err != nil {
			return rating, err
		}
	} else {
		return rating, errors.New(fmt.Sprintf("failed to get rating for play %d", date))
	}

	return rating, nil
}

func validatePlay(play PlayInfo) error {
	// check note counts add up to total

	totalCriticalPerfect := play.TapCriticalPerfect + play.HoldCriticalPerfect +
				play.SlideCriticalPerfect + play.TouchCriticalPerfect +
				play.BreakCriticalPerfect
	if totalCriticalPerfect != play.TotalCriticalPerfect && totalCriticalPerfect != 0 {
		   return errors.New(fmt.Sprintf("error validating PlayInfo %d: no. of Critical Perfects does not add up", play.UserPlayDate))
	}

	totalPerfect := play.TapPerfect + play.HoldPerfect +
			play.SlidePerfect + play.TouchPerfect +
			play.BreakPerfect
	if totalPerfect != play.TotalPerfect && totalPerfect != 0 {
		   return errors.New(fmt.Sprintf("error validating PlayInfo %d: no. of Perfects does not add up", play.UserPlayDate))
	}

	totalGreat := play.TapGreat + play.HoldGreat +
			play.SlideGreat + play.TouchGreat +
			play.BreakGreat
	if totalGreat != play.TotalGreat && totalGreat != 0 {
		   return errors.New(fmt.Sprintf("error validating PlayInfo %d: no. of Greats does not add up", play.UserPlayDate))
	}

	totalGood := play.TapGood + play.HoldGood +
			play.SlideGood + play.TouchGood +
			play.BreakGood
	if totalGood != play.TotalGood && totalGood != 0 {
		   return errors.New(fmt.Sprintf("error validating PlayInfo %d: no. of Goods does not add up", play.UserPlayDate))
	}

	totalMiss := play.TapMiss + play.HoldMiss +
			play.SlideMiss + play.TouchMiss +
			play.BreakMiss
	if totalMiss != play.TotalMiss && totalMiss != 0 {
		   return errors.New(fmt.Sprintf("error validating PlayInfo %d: no. of Misses does not add up", play.UserPlayDate))
	}

	// TODO: more validation on combo status, score, etc.

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

// GetPlay returns a PlayInfo that corresponds to date
func (playdb *PlayDB) GetPlay(date int64) (PlayInfo, error) {
	rows, err := playdb.db.Query(`
	SELECT * FROM plays WHERE user_play_date=?`, date)
	if err != nil {
		return PlayInfo{}, err
	}
	defer rows.Close()

	plays, err := rowsToPlayInfos(rows)
	if err != nil {
		return PlayInfo{}, err
	}

	if len(plays) < 1 {
		return PlayInfo{}, &PlayNotFoundError{UserPlayDate: date}
	}

	return plays[0], nil
}

// GetPlays returns a slice of PlayInfos.
// ascending: whether dates are ascending or descending.
// limit: the maximum length of the slice.
// offset: offset in the database.
func (playdb *PlayDB) GetPlays(ascending bool, limit, offset int) ([]PlayInfo, error) {
	var stmt *sql.Stmt
	if ascending {
		stmt = playdb.ascStmt
	} else {
		stmt = playdb.descStmt
	}

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plays, err := rowsToPlayInfos(rows)
	if err != nil {
		return nil, err
	}

	return plays, nil
}

type PlayNotFoundError struct {
	UserPlayDate int64
}

func (e *PlayNotFoundError) Error() string {
	return fmt.Sprintf("Playlog entry with date %d not found in database", e.UserPlayDate)
}

// GetCount returns the number of plays in the database
func (playdb *PlayDB) GetCount() (int, error) {
	var count int

	rows, err := playdb.db.Query(`SELECT COUNT(*) FROM plays`)
	if err != nil {
		return count, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return count, err
		}
	} else {
		return count, errors.New("failed to get count of plays")
	}

	return count, nil
}

func (playdb *PlayDB) getCount_v1_2() (int, error) {
	var count int

	rows, err := playdb.db.Query(`
		SELECT COUNT(*) FROM dx_rating_gen_3`)
	if err != nil {
		return count, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return count, err
		}
	} else {
		return count, errors.New("failed to get count of dx_rating_gen_3")
	}

	return count, nil
}

// GetBestScoreBeforeDate returns the best score for a given song and difficulty before the specified date
func (playdb *PlayDB) GetBestScoreBeforeDate(songId int, difficulty Difficulty, date int64) (int, error) {
	var score int

	rows, err := playdb.db.Query(`
		SELECT score FROM plays WHERE song_id=? AND difficulty=? AND user_play_date<? ORDER BY score DESC LIMIT 1`,
		songId, difficulty, date)
	if err != nil {
		return score, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&score)
		if err != nil {
			return score, err
		}
	} else if err = rows.Err(); err != nil {
		return score, err
	} else {
		score = 0
	}

	return score, nil
}

// GetBestPlaysInVersion returns the n best plays in version
// GetBestPlaysNotInVersion returns the n best plays not in version

func rowsToPlayInfos(rows *sql.Rows) ([]PlayInfo, error) {
	plays := make([]PlayInfo, 0, 50)

	for rows.Next() {
		var matchingUsersJSON []byte
		play := PlayInfo{}
		err := rows.Scan(
			&play.UserPlayDate, &play.SongId, &play.Difficulty,

			&play.Score, &play.DxScore, &play.ComboStatus, &play.SyncStatus,
			&play.IsClear, &play.IsNewRecord, &play.IsDxNewRecord,
			&play.Track, &matchingUsersJSON,

			&play.MaxCombo, &play.TotalCombo, &play.MaxSync, &play.TotalSync,

			&play.FastCount, &play.LateCount, &play.BeforeRating, &play.AfterRating,

			&play.TapCriticalPerfect, &play.TapPerfect,
			&play.TapGreat, &play.TapGood, &play.TapMiss,

			&play.HoldCriticalPerfect, &play.HoldPerfect,
			&play.HoldGreat, &play.HoldGood, &play.HoldMiss,

			&play.SlideCriticalPerfect, &play.SlidePerfect,
			&play.SlideGreat, &play.SlideGood, &play.SlideMiss,

			&play.TouchCriticalPerfect, &play.TouchPerfect,
			&play.TouchGreat, &play.TouchGood, &play.TouchMiss,

			&play.BreakCriticalPerfect, &play.BreakPerfect,
			&play.BreakGreat, &play.BreakGood, &play.BreakMiss,

			&play.TotalCriticalPerfect, &play.TotalPerfect,
			&play.TotalGreat, &play.TotalGood, &play.TotalMiss)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(matchingUsersJSON, &play.MatchingUsers)
		if err != nil {
			return nil, err
		}

		plays = append(plays, play)
	}
	if err := rows.Err(); err != nil {
		return plays, err
	}

	return plays, nil
}
