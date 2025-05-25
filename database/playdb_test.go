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
	if _, ok := err.(*database.PlayNotFoundError); ok {
		t.Log("correctly returned PlayNotFoundError:", err)
	} else if err != nil {
		t.Error("returned non-nil error, but not PlayNotFoundError:", err)
	} else {
		t.Error("expected PlayNotFoundError for non-existant play")
	}

	if !reflect.DeepEqual(play1, play1g) {
		t.Log(play1)
		t.Log(play1g)
		t.Error("play1 not equal to play1g")
	}

	plays2, err := playdb.GetPlays(false, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(plays2) != 1 {
		t.Fatal("len(plays2) != 1")
	}

	if !reflect.DeepEqual(play1, plays2[0]) {
		t.Log(play1)
		t.Log(plays2[0])
		t.Error("play1 not equal to plays2[0]")
	}
}

func TestGetPlays(t *testing.T) {
	db, err := sql.Open("sqlite3", "../test/test-plays.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	playdb, err := database.NewPlayDB(db)
	if err != nil {
		t.Fatal(err)
	}

	plays1, err := playdb.GetPlays(true, 3, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(plays1) != 3 {
		t.Fatal("len(plays1) != 3")
	}

	if plays1[0].UserPlayDate != 1743108219 ||
	   plays1[1].UserPlayDate != 1743109338 ||
	   plays1[2].UserPlayDate != 1743109538 {
		t.Fatal("plays1 incorrect")
	}

	plays2, err := playdb.GetPlays(false, 4, 5)
	if err != nil {
		t.Fatal(err)
	}

	if len(plays2) != 4 {
		t.Fatal("len(plays2) != 4")
	}

	if plays2[0].UserPlayDate != 1746509521 ||
	   plays2[1].UserPlayDate != 1746509145 ||
	   plays2[2].UserPlayDate != 1746508978 ||
	   plays2[3].UserPlayDate != 1746508788 {
		t.Fatal("plays2 incorrect")
	}

	plays3, err := playdb.GetPlays(false, 5, 198)
	if err != nil {
		t.Fatal(err)
	}

	if len(plays3) != 2 {
		t.Fatal("len(plays3) != 2")
	}

	if plays3[0].UserPlayDate != 1743108219 ||
	   plays3[1].UserPlayDate != 1743108003 {
		t.Fatal("plays3 incorrect")
	}
}

func TestGetCount(t *testing.T) {
	db, err := sql.Open("sqlite3", "../test/test-plays.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	playdb, err := database.NewPlayDB(db)
	if err != nil {
		t.Fatal(err)
	}

	count, err := playdb.GetCount()
	if err != nil {
		t.Fatal(err)
	}

	if count != 200 {
		t.Error("count != 200")
	}
}

func TestGetBestScoreBeforeDate(t *testing.T) {
	db, err := sql.Open("sqlite3", "../test/test-plays.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	playdb, err := database.NewPlayDB(db)
	if err != nil {
		t.Fatal(err)
	}

	// Aibao Dance Hall
	score1, err := playdb.GetBestScoreBeforeDate(11765, database.Master, 1744485640)
	if err != nil {
		t.Fatal(err)
	}
	if score1 != 0 {
		t.Error("score1 != 0")
	}

	// Aibao Dance Hall
	score2, err := playdb.GetBestScoreBeforeDate(11765, database.Master, 1745707467)
	if err != nil {
		t.Fatal(err)
	}
	if score2 != 971931 {
		t.Error("score2 != 971931")
	}

	// Aibao Dance Hall
	score3, err := playdb.GetBestScoreBeforeDate(11765, database.Master, 1746506710)
	if err != nil {
		t.Fatal(err)
	}
	if score3 != 971931 {
		t.Error("score3 != 971931")
	}

	// Override
	score4, err := playdb.GetBestScoreBeforeDate(11794, database.Master, 1743569808)
	if err != nil {
		t.Fatal(err)
	}
	if score4 != 0 {
		t.Error("score4 != 0")
	}

	// Override
	score5, err := playdb.GetBestScoreBeforeDate(11794, database.Master, 1744401821)
	if err != nil {
		t.Fatal(err)
	}
	if score5 != 981938 {
		t.Error("score5 != 981938")
	}

	// Override
	score6, err := playdb.GetBestScoreBeforeDate(11794, database.Master, 1744922642)
	if err != nil {
		t.Fatal(err)
	}
	if score6 != 985903 {
		t.Error("score6 != 985903")
	}

	// Override
	score7, err := playdb.GetBestScoreBeforeDate(11794, database.Master, 1745701086)
	if err != nil {
		t.Fatal(err)
	}
	if score7 != 985903 {
		t.Error("score7 != 985903")
	}

	// Override
	score8, err := playdb.GetBestScoreBeforeDate(11794, database.Master, 1745701087)
	if err != nil {
		t.Fatal(err)
	}
	if score8 != 988921 {
		t.Error("score8 != 988921")
	}
}

func TestAddInvalidPlays(t *testing.T) {
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
		t.Error("non-nil error for play1")
	}

	play2 := play1
	play2.TotalCriticalPerfect += 1
	err = playdb.AddPlay(play2)
	if err == nil {
		t.Error("expected non-nil error for play2")
	} else {
		t.Log("edited crit. perfects:", err)
	}

	play3 := play1
	play3.TotalPerfect += 1
	err = playdb.AddPlay(play3)
	if err == nil {
		t.Error("expected non-nil error for play3")
	} else {
		t.Log("edited perfects:", err)
	}

	play4 := play1
	play4.TotalGreat += 1
	err = playdb.AddPlay(play4)
	if err == nil {
		t.Error("expected non-nil error for play4")
	} else {
		t.Log("edited greats:", err)
	}

	play5 := play1
	play5.TotalGood += 1
	err = playdb.AddPlay(play5)
	if err == nil {
		t.Error("expected non-nil error for play5")
	} else {
		t.Log("edited goods:", err)
	}

	play6 := play1
	play6.TotalMiss += 1
	err = playdb.AddPlay(play6)
	if err == nil {
		t.Error("expected non-nil error for play6")
	} else {
		t.Log("edited misses:", err)
	}
}

func TestAddPlayWithNoDetailedJudgement(t *testing.T) {
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

		TapCriticalPerfect	: 0,
	        TapPerfect		: 0,
		TapGreat		: 0,
		TapGood			: 0,
		TapMiss			: 0,

		HoldCriticalPerfect	: 0,
	        HoldPerfect		: 0,
		HoldGreat		: 0,
		HoldGood		: 0,
		HoldMiss		: 0,

		SlideCriticalPerfect	: 0,
	        SlidePerfect		: 0,
		SlideGreat		: 0,
		SlideGood		: 0,
		SlideMiss		: 0,

		TouchCriticalPerfect	: 0,
	        TouchPerfect		: 0,
		TouchGreat		: 0,
		TouchGood		: 0,
		TouchMiss		: 0,

		BreakCriticalPerfect	: 0,
	        BreakPerfect		: 0,
		BreakGreat		: 0,
		BreakGood		: 0,
		BreakMiss		: 0,

		TotalCriticalPerfect	: 393,
	        TotalPerfect		: 290,
		TotalGreat		: 82,
		TotalGood		: 13,
		TotalMiss		: 5,
	}

	err = playdb.AddPlay(play1)
	if err != nil {
		t.Error("non-nil error for play1")
	}
}
