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

// Package solips handles making api requests to solips.app and
// updating the play database
package solips

import (
	"fmt"
	"log"
	"io"
	"time"
	"strings"
	"net/http"
	"net/http/cookiejar"
	"encoding/json"
	"errors"

	"github.com/yadayadajaychan/playlog/database"
	"github.com/yadayadajaychan/playlog/internal/context"
)

const (
	apiUrl = `https://www.solips.app/api/trpc/maimai.playlog,maimai.favorites?batch=1&input={"0":{"json":null,"meta":{"values":["undefined"]}},"1":{"json":null,"meta":{"values":["undefined"]}}}`
	loginUrl = `https://www.solips.app/api/trpc/card.link?batch=1`
	detailUrlFmt = `https://www.solips.app/api/trpc/maimai.playlogDetail,maimai.favorites?batch=1&input={"0":{"json":{"playlogId":"%s"}},"1":{"json":null,"meta":{"values":["undefined"]}}}`
	playlogLength = 100
)

type apiPlaylogV2 struct {
	Result	struct {
		Data	struct {
			Json	[]apiPlaylogEntry
		}
	}
}

type apiPlaylog struct {
	Playlog	[]apiPlaylogEntry
}

type apiPlaylogEntry struct {
	PlaylogApiId	string
	Info		struct {
		UserPlayDate	string
	}
}

type apiPlaylogDetailV2 struct {
	Result	struct {
		Data	struct {
			JSON	maimaiPlaylogDetail
		}
	}
}

type apiPlaylogDetail struct {
	MaimaiPlaylogDetail	maimaiPlaylogDetail
}

type maimaiPlaylogDetail struct {
	Info	struct {
		MusicId		int
		Level		string
		Achievement	int
		Deluxscore	int
		ScoreRank	string
		ComboStatus	string
		SyncStatus	string
		IsClear		bool
		IsAchieveNewRecord	bool
		IsDeluxscoreNewRecord	bool
		Track		int
		UserPlayDate	string
	}

	Detail	struct {
		JudgeTap	struct {
			TapCriticalPerfect	int
			TapPerfect		int
			TapGreat		int
			TapGood			int
			TapMiss			int
		}
		JudgeHold	struct {
			HoldCriticalPerfect	int
			HoldPerfect		int
			HoldGreat		int
			HoldGood		int
			HoldMiss		int
		}
		JudgeSlide	struct {
			SlideCriticalPerfect	int
			SlidePerfect		int
			SlideGreat		int
			SlideGood		int
			SlideMiss		int
		}
		JudgeTouch	struct {
			TouchCriticalPerfect	int
			TouchPerfect		int
			TouchGreat		int
			TouchGood		int
			TouchMiss		int
		}
		JudgeBreak	struct {
			BreakCriticalPerfect	int
			BreakPerfect		int
			BreakGreat		int
			BreakGood		int
			BreakMiss		int
		}

		MaxCombo	int
		TotalCombo	int
		MaxSync		int
		TotalSync	int
		FastCount	int
		LateCount	int
		BeforeRating	int
		AfterRating	int
	}

	MatchingUsers []struct {
		UserName	string
	}
}

var globalCookieJar, _ = cookiejar.New(nil)

// Update uses the Mythos access code to get the most recent 100 songs played
// and makes an api request per new song that's not in the database,
// delaying by ctx.ApiInterval between requests.
// It then adds them to the database.
// ctx requires Playdb, Songdb, AccessCode, ApiInterval, Verbose
func Update(ctx context.PlaylogCtx) error {
	playlog, err := getPlaylog(ctx.AccessCode)
	if err != nil {
		return err
	}

	err = validatePlaylog(playlog)
	if err != nil {
		return err
	}

	for _, entry := range playlog.Playlog {
		playdate, err := time.Parse(time.RFC3339, entry.Info.UserPlayDate)
		if err != nil {
			return err
		}

		_, err = ctx.Playdb.GetPlay(playdate.Unix())
		if _, ok := err.(*database.PlayNotFoundError); ok {
			playlogDetail, err := getPlaylogDetail(ctx.AccessCode, entry.PlaylogApiId)
			if err != nil {
				return err
			}

			err = validatePlaylogDetail(playlogDetail, ctx.Songdb)
			if err != nil {
				return err
			}

			err = addMaimaiPlaylogDetailToPlayDB(ctx.Playdb, playlogDetail.MaimaiPlaylogDetail)
			if err != nil {
				return err
			}

			if ctx.Verbose >= 1 {
				log.Printf("play %d: added to db\n", playdate.Unix())
			}
			time.Sleep(ctx.ApiInterval)

		} else if err != nil {
			return err
		} else {
			if ctx.Verbose >= 2 {
				log.Printf("play %d: already exists in db\n", playdate.Unix())
			}
		}
	}


	return nil
}

func levelToDifficulty(lvl string) (database.Difficulty, error) {
	var difficulty database.Difficulty
	switch lvl {
	case "MAIMAI_LEVEL_BASIC":
		difficulty = database.Basic
	case "MAIMAI_LEVEL_ADVANCED":
		difficulty = database.Advanced
	case "MAIMAI_LEVEL_EXPERT":
		difficulty = database.Expert
	case "MAIMAI_LEVEL_MASTER":
		difficulty = database.Master
	case "MAIMAI_LEVEL_REMASTER":
		difficulty = database.ReMaster
	case "MAIMAI_LEVEL_UTAGE":
		difficulty = database.Utage
	default:
		return difficulty, errors.New("invalid level: " + lvl)
	}

	return difficulty, nil
}

func addMaimaiPlaylogDetailToPlayDB(playdb *database.PlayDB, maimai maimaiPlaylogDetail) error {
	playdate, err := time.Parse(time.RFC3339, maimai.Info.UserPlayDate)
	if err != nil {
		return err
	}

	difficulty, err := levelToDifficulty(maimai.Info.Level)
	if err != nil {
		return err
	}

	var comboStatus database.ComboStatus
	switch maimai.Info.ComboStatus {
	case "MAIMAI_COMBO_STATUS_NONE":
		comboStatus = database.NoCombo
	case "MAIMAI_COMBO_STATUS_FULL_COMBO":
		comboStatus = database.FullCombo
	case "MAIMAI_COMBO_STATUS_FULL_COMBO_PLUS":
		comboStatus = database.FullComboPlus
	case "MAIMAI_COMBO_STATUS_ALL_PERFECT":
		comboStatus = database.AllPerfect
	case "MAIMAI_COMBO_STATUS_ALL_PERFECT_PLUS":
		comboStatus = database.AllPerfectPlus
	default:
		return errors.New("addMaimaiPlaylogDetailToPlayDB: invalid combo status: " + maimai.Info.ComboStatus)
	}

	var syncStatus database.SyncStatus
	switch maimai.Info.SyncStatus {
	case "MAIMAI_SYNC_STATUS_NONE":
		syncStatus = database.NoSync
	case "MAIMAI_SYNC_STATUS_FULL_SYNC":
		syncStatus = database.FullSync
	case "MAIMAI_SYNC_STATUS_FULL_SYNC_PLUS":
		syncStatus = database.FullSyncPlus
	case "MAIMAI_SYNC_STATUS_FULL_SYNC_DX":
		syncStatus = database.FullSyncDx
	case "MAIMAI_SYNC_STATUS_FULL_SYNC_DX_PLUS":
		syncStatus = database.FullSyncDxPlus
	default:
		return errors.New("addMaimaiPlaylogDetailToPlayDB: invalid sync status: " + maimai.Info.SyncStatus)
	}

	matchingUsers := make([]string, 0, 1)
	for _, v := range maimai.MatchingUsers {
		matchingUsers = append(matchingUsers, v.UserName)
	}

	playinfo := database.PlayInfo{
		UserPlayDate : playdate.Unix(),
		SongId       : maimai.Info.MusicId,
		Difficulty   : difficulty,

		Score         : maimai.Info.Achievement,
		DxScore       : maimai.Info.Deluxscore,
		ComboStatus   : comboStatus,
		SyncStatus    : syncStatus,
		IsClear       : maimai.Info.IsClear,
		IsNewRecord   : maimai.Info.IsAchieveNewRecord,
		IsDxNewRecord : maimai.Info.IsDeluxscoreNewRecord,
		Track         : maimai.Info.Track,
		MatchingUsers : matchingUsers,

		MaxCombo   : maimai.Detail.MaxCombo,
		TotalCombo : maimai.Detail.TotalCombo,
		MaxSync    : maimai.Detail.MaxSync,
		TotalSync  : maimai.Detail.TotalSync,

		FastCount    : maimai.Detail.FastCount,
		LateCount    : maimai.Detail.LateCount,
		BeforeRating : maimai.Detail.BeforeRating,
		AfterRating  : maimai.Detail.AfterRating,

		TapCriticalPerfect : maimai.Detail.JudgeTap.TapCriticalPerfect,
		TapPerfect         : maimai.Detail.JudgeTap.TapPerfect,
		TapGreat           : maimai.Detail.JudgeTap.TapGreat,
		TapGood            : maimai.Detail.JudgeTap.TapGood,
		TapMiss            : maimai.Detail.JudgeTap.TapMiss,

		HoldCriticalPerfect : maimai.Detail.JudgeHold.HoldCriticalPerfect,
		HoldPerfect         : maimai.Detail.JudgeHold.HoldPerfect,
		HoldGreat           : maimai.Detail.JudgeHold.HoldGreat,
		HoldGood            : maimai.Detail.JudgeHold.HoldGood,
		HoldMiss            : maimai.Detail.JudgeHold.HoldMiss,

		SlideCriticalPerfect : maimai.Detail.JudgeSlide.SlideCriticalPerfect,
		SlidePerfect         : maimai.Detail.JudgeSlide.SlidePerfect,
		SlideGreat           : maimai.Detail.JudgeSlide.SlideGreat,
		SlideGood            : maimai.Detail.JudgeSlide.SlideGood,
		SlideMiss            : maimai.Detail.JudgeSlide.SlideMiss,

		TouchCriticalPerfect : maimai.Detail.JudgeTouch.TouchCriticalPerfect,
		TouchPerfect         : maimai.Detail.JudgeTouch.TouchPerfect,
		TouchGreat           : maimai.Detail.JudgeTouch.TouchGreat,
		TouchGood            : maimai.Detail.JudgeTouch.TouchGood,
		TouchMiss            : maimai.Detail.JudgeTouch.TouchMiss,

		BreakCriticalPerfect : maimai.Detail.JudgeBreak.BreakCriticalPerfect,
		BreakPerfect         : maimai.Detail.JudgeBreak.BreakPerfect,
		BreakGreat           : maimai.Detail.JudgeBreak.BreakGreat,
		BreakGood            : maimai.Detail.JudgeBreak.BreakGood,
		BreakMiss            : maimai.Detail.JudgeBreak.BreakMiss,

		TotalCriticalPerfect : maimai.Detail.JudgeTap.TapCriticalPerfect +
					maimai.Detail.JudgeHold.HoldCriticalPerfect +
					maimai.Detail.JudgeSlide.SlideCriticalPerfect +
					maimai.Detail.JudgeTouch.TouchCriticalPerfect +
					maimai.Detail.JudgeBreak.BreakCriticalPerfect,

		TotalPerfect : maimai.Detail.JudgeTap.TapPerfect +
				maimai.Detail.JudgeHold.HoldPerfect +
				maimai.Detail.JudgeSlide.SlidePerfect +
				maimai.Detail.JudgeTouch.TouchPerfect +
				maimai.Detail.JudgeBreak.BreakPerfect,

		TotalGreat : maimai.Detail.JudgeTap.TapGreat +
				maimai.Detail.JudgeHold.HoldGreat +
				maimai.Detail.JudgeSlide.SlideGreat +
				maimai.Detail.JudgeTouch.TouchGreat +
				maimai.Detail.JudgeBreak.BreakGreat,

		TotalGood : maimai.Detail.JudgeTap.TapGood +
				maimai.Detail.JudgeHold.HoldGood +
				maimai.Detail.JudgeSlide.SlideGood +
				maimai.Detail.JudgeTouch.TouchGood +
				maimai.Detail.JudgeBreak.BreakGood,

		TotalMiss : maimai.Detail.JudgeTap.TapMiss +
				maimai.Detail.JudgeHold.HoldMiss +
				maimai.Detail.JudgeSlide.SlideMiss +
				maimai.Detail.JudgeTouch.TouchMiss +
				maimai.Detail.JudgeBreak.BreakMiss,
	}

	err = playdb.AddPlay(playinfo)
	if err != nil {
		return err
	}

	return nil
}

// getPlaylog gets the non-detailed playlog of the most recent 100 plays.
// only the playlogApiId and userPlayDate values matter in this case.
func getPlaylog(accessCode string) (*apiPlaylog, error) {
	// POST to loginUrl and save cookie
	client := &http.Client{
		Jar: globalCookieJar,
	}

	loginData := fmt.Sprintf(`{"0":{"json":"%s"}}`, accessCode)
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(loginData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// GET with authentication cookie
	req, err = http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal JSON
	var playlogV2 apiPlaylogV2
	dec := json.NewDecoder(resp.Body)

	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '[' {
		return nil, fmt.Errorf("expected [")
	}

	if err = dec.Decode(&playlogV2); err != nil {
		return nil, err
	}

	// convert from V2 to V1
	playlog := &apiPlaylog{
		Playlog: playlogV2.Result.Data.Json,
	}

	return playlog, nil
}

func validatePlaylog(playlog *apiPlaylog) error {
	// verify there are playlogLength items
	n := len(playlog.Playlog)
	if n != playlogLength {
		return errors.New(fmt.Sprintf("len(playlog): expected %d, got %d", playlogLength, n))
	}

	// check for duplicates
	seenPlaylogApiId := make(map[string]bool)
	seenUserPlayDate := make(map[string]bool)
	for _, item := range playlog.Playlog {
		if seenPlaylogApiId[item.PlaylogApiId] {
			return errors.New(fmt.Sprint("duplicate PlaylogApiId: ", item.PlaylogApiId))
		}
		seenPlaylogApiId[item.PlaylogApiId] = true

		if seenUserPlayDate[item.Info.UserPlayDate] {
			return errors.New(fmt.Sprint("duplicate UserPlayDate: ", item.Info.UserPlayDate))
		}
		seenUserPlayDate[item.Info.UserPlayDate] = true
	}

	// TODO: validate dates

	return nil
}

// convenience function for diagnostics
func printPlaylog(playlog *apiPlaylog) {
	for _, item := range playlog.Playlog {
		fmt.Printf("%v\t%v\n", item.PlaylogApiId, item.Info.UserPlayDate)
	}
}

// getPlaylogDetail depends on globalCookieJar to work correctly.
func getPlaylogDetail(accessCode, playlogApiId string) (*apiPlaylogDetail, error) {
	client := &http.Client{
		Jar: globalCookieJar,
	}

	url := fmt.Sprintf(detailUrlFmt, playlogApiId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal JSON
	var playlogDetailV2 apiPlaylogDetailV2
	dec := json.NewDecoder(resp.Body)

	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '[' {
		return nil, fmt.Errorf("expected [")
	}

	if err = dec.Decode(&playlogDetailV2); err != nil {
		return nil, err
	}

	// convert from V2 to V1
	playlogDetail := &apiPlaylogDetail{
		MaimaiPlaylogDetail: playlogDetailV2.Result.Data.JSON,
	}

	return playlogDetail, nil
}

func validatePlaylogDetail(playlogDetail *apiPlaylogDetail, songdb *database.SongDB) error {
	detail := playlogDetail.MaimaiPlaylogDetail

	if len(detail.Info.UserPlayDate) <= 0 {
		return errors.New("error validating playlogDetail, userPlayDate not found")
	}

	song, err := songdb.GetSong(detail.Info.MusicId)
	if err != nil {
		return err
	}

	difficulty, err := levelToDifficulty(detail.Info.Level)
	if err != nil {
		return err
	}

	totalCombo := detail.Detail.TotalCombo
	if totalCombo == 0 {
		// TODO add up all the notes instead
		return nil
	}

	for _, chart := range song.Charts {
		if chart.Difficulty == difficulty {
			if chart.MaxNotes == totalCombo || chart.MaxNotes == 0 {
				return nil
			}
		}
	}

	return errors.New(fmt.Sprintf("error validating playlogDetail, invalid songId: %d", detail.Info.MusicId))
}
