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
	//"github.com/yadayadajaychan/playlog/database"
	"github.com/yadayadajaychan/playlog/internal/context"
)

const (
	apiUrl = "https://kamai.tachi.ac/api/v1"
)

func Update(ctx context.PlaylogCtx) error {
	return nil
}

type activityJSON struct {
	Success bool
	Body struct {
		RecentSessions []sessionJSON
	}
}

type sessionJSON struct {
	ScoreIDs    []string
	TimeStarted int
}

type scoreJSON struct {
	Success bool
	Body struct {
		Score struct {
			TimeAchieved int // unix time in milliseconds
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
