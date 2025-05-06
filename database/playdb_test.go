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

package database_test

import (
	"testing"
	"os"
	"reflect"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yadayadajaychan/playlog/database"
)

func TestAddAndGetPlay(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "playdb-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	db, err := sql.Open("sqlite3", tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	playdb, err := database.NewPlayDB(db)
	if err != nil {
		t.Fatal(err)
	}

	play1 := database.PlayInfo{
		UserPlayDate	: 1743108003,
		SongId		: 11441,
		Difficulty	: database.Master,

		Score		: 971017,
		DxScore		: 1841,
		ComboStatus	: database.NoCombo,
		SyncStatus	: database.NoSync,
		IsClear		: true,
		IsNewRecord	: true,
		IsDxNewRecord	: true,
		Track		: 3,
		MatchingUsers	: []string{"ＳＵＰＡＩＤＯＬ"},

		MaxCombo	: 385,
		TotalCombo	: 783,
		MaxSync		: 559,
		TotalSync	: 1566,

		FastCount	: 53,
		LateCount	: 66,
		BeforeRating	: 13085,
		AfterRating	: 13085,

		TapCriticalPerfect	: 222,
	        TapPerfect		: 239,
		TapGreat		: 67,
		TapGood			: 8,
		TapMiss			: 3,

		HoldCriticalPerfect	: 44,
	        HoldPerfect		: 27,
		HoldGreat		: 6,
		HoldGood		: 1,
		HoldMiss		: 1,

		SlideCriticalPerfect	: 93,
	        SlidePerfect		: 0,
		SlideGreat		: 3,
		SlideGood		: 3,
		SlideMiss		: 0,

		TouchCriticalPerfect	: 19,
	        TouchPerfect		: 0,
		TouchGreat		: 0,
		TouchGood		: 0,
		TouchMiss		: 1,

		BreakCriticalPerfect	: 15,
	        BreakPerfect		: 24,
		BreakGreat		: 6,
		BreakGood		: 1,
		BreakMiss		: 0,

		TotalCriticalPerfect	: 393,
	        TotalPerfect		: 290,
		TotalGreat		: 82,
		TotalGood		: 13,
		TotalMiss		: 5,
	}

	err = playdb.AddPlay(play1)
	if err != nil {
		t.Fatal(err)
	}

	play1g, err := playdb.GetPlay(1743108003)
	if err != nil {
		t.Fatal(err)
	}

	_, err = playdb.GetPlay(123)
	if err != nil {
		t.Log("correctly returned non-nil error:", err)
	} else {
		t.Error("expected non-nil error for non-existant play")
	}

	if !reflect.DeepEqual(play1, play1g) {
		t.Log(play1)
		t.Log(play1g)
		t.Error("play1 not equal")
	}
}
