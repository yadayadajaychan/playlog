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

// Package update handles making api requests to solips.app and
// updating the playlog database
package update

import (
	"fmt"
	"io"
	"time"
	"strings"
	"net/url"
	"net/http"
	"net/http/cookiejar"
	"encoding/json"
	"errors"
)

const (
	apiUrl = "https://www.solips.app/maimai/profile?_data=routes%2Fmaimai.profile"
	playlogLength = 100
)

type apiPlaylog struct {
	Playlog	[]apiPlaylogEntry
}

type apiPlaylogEntry struct {
	PlaylogApiId	string
	Info		struct {
		UserPlayDate	string
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

// Update uses the Mythos access code to get the most recent 100 songs played
// and makes an api request per new song that's not in the database,
// delaying by apiDelay between requests. It then adds them to the database.
func Update(accessCode string, apiDelay time.Duration) error {
	playlog, err := getPlaylog(accessCode)
	if err != nil {
		return err
	}

	err = validatePlaylog(playlog)
	if err != nil {
		return err
	}

	printPlaylog(playlog)

	//time.Sleep(apiDelay)

	return nil
}

// getPlaylog gets the non-detailed playlog of the most recent 100 plays.
// only the playlogApiId and userPlayDate values matter in this case.
func getPlaylog(accessCode string) (*apiPlaylog, error) {
	// POST to API and save cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	form := url.Values{
		"accessCode": {accessCode},
		"requestType": {"getUserApiId"},
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// GET with authentication cookie
	req.Method = "GET"
	req.Body = nil

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Unmarshal JSON
	playlog := &apiPlaylog{
		Playlog: make([]apiPlaylogEntry, 0, 100),
	}
	err = json.Unmarshal(data, playlog)
	if err != nil {
		return nil, err
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

func getPlaylogDetail(accessCode, playlogApiId string) (*apiPlaylogDetail, error) {
	form := url.Values{
		"accessCode": {accessCode},
		"requestType": {"getPlaylogDetail"},
		"playlogApiId": {playlogApiId},
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	playlogDetail := &apiPlaylogDetail{}

	err = json.Unmarshal(data, playlogDetail)
	if err != nil {
		return nil, err
	}

	return playlogDetail, nil
}

func validatePlaylogDetail(playlogDetail *apiPlaylogDetail) error {
	if len(playlogDetail.MaimaiPlaylogDetail.Info.UserPlayDate) <= 0 {
		return errors.New("error validating playlogDetail, userPlayDate not found")
	}

	return nil
}
