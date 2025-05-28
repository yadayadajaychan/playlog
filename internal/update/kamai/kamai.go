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

// Package kamai handles making api requests to kamai.tachi.ac and
// updating the play database
package kamai

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"encoding/json"
	"time"
	"strings"
	"math"
	"errors"

	"github.com/yadayadajaychan/playlog/database"
	"github.com/yadayadajaychan/playlog/internal/context"
)

const (
	apiUrl = "https://kamai.tachi.ac/api/v1"
)

func Update(ctx context.PlaylogCtx) error {
	sess := &sessions{User: ctx.KamaiUser}
	allScoreIds := make([]string, 0, 100)

	for sess.Next() {
		ss := sess.Get()
		scoreIds := make([]string, 0, 30)
		for _, s := range ss {
			scoreIds = append(scoreIds, s.ScoreIDs...)
		}
		allScoreIds = append(allScoreIds, scoreIds...)

		if ctx.Verbose >= 1 {
			log.Printf("retrieved %d scoreIds", len(scoreIds))
		}

		time.Sleep(ctx.ApiInterval)
	}
	if sess.Err() != nil {
		return sess.Err()
	}

	for _, scoreId := range allScoreIds {
		score, err := getScore(scoreId)
		if err != nil {
			return err
		}

		err = addScoreToPlayDB(score, ctx)
		if err != nil {
			return err
		}

		if ctx.Verbose >= 1 {
			log.Printf("added play %d to database", score.Body.Score.TimeAchieved / 1000)
		}
		time.Sleep(ctx.ApiInterval)
	}

	return nil
}

func kamaiDiffToDiff(kamaiDifficulty string) (database.Difficulty, error) {
	kamaiDifficulty = strings.ToLower(kamaiDifficulty)
	var difficulty database.Difficulty
	switch kamaiDifficulty {
	case "basic", "dx basic":
		difficulty = database.Basic
	case "advanced", "dx advanced":
		difficulty = database.Advanced
	case "expert", "dx expert":
		difficulty = database.Expert
	case "master", "dx master":
		difficulty = database.Master
	case "re:master", "dx re:master":
		difficulty = database.ReMaster
	default:
		return difficulty, errors.New("invalid difficulty: " + kamaiDifficulty)
	}

	return difficulty, nil
}

func toSongType(kamaiDifficulty string) string {
	kamaiDifficulty = strings.ToLower(kamaiDifficulty)
	if kamaiDifficulty[0:2] == "dx" {
		return "dx"
	}

	return "std"
}

func addScoreToPlayDB(score scoreJSON, ctx context.PlaylogCtx) error {
	playDate := score.Body.Score.TimeAchieved / 1000

	_, err := ctx.Playdb.GetPlay(playDate)
	if err == nil {
		if ctx.Verbose >= 2 {
			log.Printf("play %d already exists in db\n", playDate)
		}
		return nil
	} else if _, ok := err.(*database.PlayNotFoundError); !ok {
		return err
	}

	scoreData := score.Body.Score.ScoreData

	songs, err := ctx.Songdb.GetSongsByName(score.Body.Song.Title)
	if err != nil {
		return err
	}

	difficulty, err := kamaiDiffToDiff(score.Body.Chart.Difficulty)
	if err != nil {
		return err
	}

	songType := toSongType(score.Body.Chart.Difficulty)

	var song *database.SongInfo
	for _, s := range songs {
		if songType == s.Type {
			song = &s
		}
	}
	if song == nil {
		return errors.New(fmt.Sprintf("no song with name '%s' and type '%s' found", score.Body.Song.Title, songType))
	}

	var chart *database.ChartInfo
	for _, c := range song.Charts {
		if c.Difficulty == difficulty {
			chart = &c
		}
	}
	if chart == nil {
		return errors.New(fmt.Sprintf("no chart with difficulty '%d' for song with name '%s' and type '%s' found", difficulty, score.Body.Song.Title, songType))
	}

	kamaiLevel := int(math.Round(score.Body.Chart.LevelNum * 10))
	if kamaiLevel != chart.InternalLevel {
		return errors.New(fmt.Sprintf("kamai level '%d' does not match internal level '%d' for song '%s'", kamaiLevel, chart.InternalLevel, score.Body.Song.Title))
	}

	play := database.PlayInfo{
		UserPlayDate : playDate,
		SongId       : song.SongId,
		Difficulty   : difficulty,

		Score: int(math.Round(scoreData.Percent * 10000)),
		//DxScore: 0,
		//ComboStatus: ,
		//SyncStatus: 0,
	}
	_=play

	return nil
}

func getScore(scoreId string) (scoreJSON, error) {
	score := scoreJSON{}
	url := apiUrl + "/scores/" + scoreId + "?getRelated"

	resp, err := http.Get(url)
	if err != nil {
		return score, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return score, err
	}
	resp.Body.Close()

	err = json.Unmarshal(data, &score)
	if err != nil {
		return score, err
	}

	return score, nil
}

type sessions struct {
	User      string
	startTime int64
	sessions  []sessionJSON
	err       error
}

// false if error or no more sessions.
// Use Err to differentiate.
func (s *sessions) Next() bool {
	url := apiUrl + "/users/" + s.User + "/games/maimaidx/Single/activity"

	var resp *http.Response
	if s.startTime == 0 {
		resp, s.err = http.Get(url)
	} else {
		resp, s.err = http.Get(url + fmt.Sprintf("?startTime=%d",s.startTime))
	}
	if s.err != nil {
		return false
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		s.err = err
		return false
	}
	resp.Body.Close()

	var activity activityJSON
	s.err = json.Unmarshal(data, &activity)
	if s.err != nil {
		return false
	}

	if !activity.Success {
		s.err = errors.New("kamai: call to activity api failed")
		return false
	}

	sessions := activity.Body.RecentSessions

	if len(sessions) <= 0 {
		return false
	}

	s.startTime = sessions[len(sessions)-1].TimeStarted
	s.sessions = sessions

	return true
}

func (s *sessions) Get() []sessionJSON {
	return s.sessions
}

func (s *sessions) Err() error {
	return s.err
}

type activityJSON struct {
	Success bool
	Body struct {
		RecentSessions []sessionJSON
	}
}

type sessionJSON struct {
	ScoreIDs    []string
	TimeStarted int64
}

type scoreJSON struct {
	Success bool
	Body struct {
		Score struct {
			TimeAchieved int64 // unix time in milliseconds
			ScoreData scoreDataJSON
		}
		Song  songDataJSON
		Chart chartDataJSON
	}
}

type scoreDataJSON struct {
	Percent float64
	Lamp    string
	Judgements struct {
		Pcrit   int
		Perfect int
		Great   int
		Good    int
		Miss    int
	}
	Optional struct {
		Fast int
		Slow int
		MaxCombo int
	}
	EnumIndexes struct {
		Lamp  int
		Grade int
	}
}

type songDataJSON struct {
	Title string
}

type chartDataJSON struct {
	Difficulty string
	LevelNum   float64
}
